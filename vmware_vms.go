package k3s

import (
	"fmt"
	"github.com/pulumi/pulumi-vsphere/sdk/v3/go/vsphere"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var vsphereIDs *vSphereIdData

type vSphereIdData struct {
	resourcePoolID string
	dataCenterID   string
	dataStoreID    string
	hostID         string
}

func setNetworkId(ctx *pulumi.Context, network *Network) error {
	n, err := vsphere.GetNetwork(ctx, &vsphere.GetNetworkArgs{
		DatacenterId: &vsphereIDs.dataCenterID,
		Name:         network.NetworkName,
	})

	if err != nil {
		return err
	}

	network.networkID = n.Id
	return nil
}

func newNetworkArray(nets []*Network) vsphere.VirtualMachineNetworkInterfaceArray {
	var networkArray vsphere.VirtualMachineNetworkInterfaceArray
	for _, n := range nets {
		network := &vsphere.VirtualMachineNetworkInterfaceArgs{
			AdapterType:  pulumi.StringPtr(n.AdapterType),
			NetworkId:    pulumi.String(n.networkID),
			UseStaticMac: pulumi.BoolPtr(n.StaticMac),
		}

		if n.StaticMac {
			network.MacAddress = pulumi.StringPtr(n.MacAddress)
		}

		networkArray = append(networkArray, network)
	}

	return networkArray
}

func newDiskArray(disks []Disk, namePrefix string) vsphere.VirtualMachineDiskArray {
	var diskArray vsphere.VirtualMachineDiskArray
	for i, d := range disks {
		diskArgs := &vsphere.VirtualMachineDiskArgs{
			DatastoreId:     pulumi.StringPtr(vsphereIDs.dataStoreID),
			EagerlyScrub:    pulumi.BoolPtr(d.EagerlyScrub),
			Label:           pulumi.StringPtr(fmt.Sprintf("%s-disk-%d.vmdk", namePrefix, i)),
			Size:            d.Size,
			ThinProvisioned: pulumi.BoolPtr(d.ThinProvisioned),
			UnitNumber:      pulumi.IntPtr(i),
		}
		diskArray = append(diskArray, diskArgs)
	}
	return diskArray
}

func newCloneArgs(templateID string) vsphere.VirtualMachineCloneArgs {
	return vsphere.VirtualMachineCloneArgs{
		Customize:    nil,
		LinkedClone:  pulumi.Bool(false),
		TemplateUuid: pulumi.String(templateID),
	}
}

func setVMOptions(conf *VMConfig) *vsphere.VirtualMachineArgs {
	var diskArray vsphere.VirtualMachineDiskArray
	var networkArray vsphere.VirtualMachineNetworkInterfaceArray
	if conf.Disks != nil {
		diskArray = newDiskArray(conf.Disks, conf.Name)
	}

	if conf.Networks != nil {
		networkArray = newNetworkArray(conf.Networks)
	}

	return &vsphere.VirtualMachineArgs{
		WaitForGuestIpTimeout: pulumi.Int(200),
		CpuHotAddEnabled:      pulumi.Bool(true),
		CpuHotRemoveEnabled:   pulumi.Bool(true),
		EnableDiskUuid:        pulumi.BoolPtr(conf.EnableDiskUuid),
		Folder:                pulumi.StringPtr(conf.Folder),
		Memory:                pulumi.IntPtr(conf.Memory),
		MemoryHotAddEnabled:   pulumi.Bool(true),
		Name:                  pulumi.StringPtr(conf.Name),
		NumCoresPerSocket:     pulumi.IntPtr(conf.Cores),
		NumCpus:               pulumi.IntPtr(conf.NumCpus),
		SyncTimeWithHost:      pulumi.Bool(true),
		ResourcePoolId:        pulumi.String(vsphereIDs.resourcePoolID),
		Disks:                 diskArray,
		NetworkInterfaces:     networkArray,
	}
}

func getDataCenterID(ctx *pulumi.Context, dataCenter *string) error {
	dc, err := vsphere.LookupDatacenter(ctx, &vsphere.LookupDatacenterArgs{Name: dataCenter})
	if err != nil {
		return err
	}
	vsphereIDs.dataCenterID = dc.Id
	return nil
}

func setResourcePoolID(ctx *pulumi.Context, poolName, clusterName string) error {
	if poolName == "" {
		cp, err := vsphere.LookupComputeCluster(ctx, &vsphere.LookupComputeClusterArgs{
			DatacenterId: &vsphereIDs.dataCenterID,
			Name:         clusterName,
		})

		if err != nil {
			return err
		}

		vsphereIDs.resourcePoolID = cp.ResourcePoolId
	} else {
		rp, err := vsphere.LookupResourcePool(ctx, &vsphere.LookupResourcePoolArgs{
			DatacenterId: &vsphereIDs.dataCenterID,
			Name:         &poolName,
		})

		if err != nil {
			return err
		}
		vsphereIDs.resourcePoolID = rp.Id
	}

	return nil
}

// Sets all resource ID's we need for building a VM
func setResourceIDS(ctx *pulumi.Context, conf *VMConfig) error {

	if err := getDataCenterID(ctx, conf.Datacenter); err != nil {
		return err
	}

	// If ClusterDatastore is passed get that ID else get Datastore ID
	if conf.ClusterDatastore != "" {
		ds, err := vsphere.LookupDatastoreCluster(ctx, &vsphere.LookupDatastoreClusterArgs{
			DatacenterId: &vsphereIDs.dataCenterID,
			Name:         conf.ClusterDatastore,
		})

		if err != nil {
			return err
		}

		vsphereIDs.dataStoreID = ds.Id
	} else {
		ds, err := vsphere.GetDatastore(ctx, &vsphere.GetDatastoreArgs{
			DatacenterId: &vsphereIDs.dataCenterID,
			Name:         string(conf.DataStore),
		})

		if err != nil {
			return err
		}
		vsphereIDs.dataStoreID = ds.Id
	}

	if err := setResourcePoolID(ctx, conf.ResourcePool, conf.Cluster); err != nil {
		return err
	}

	networks := conf.Networks

	for _, n := range networks {
		if err := setNetworkId(ctx, n); err != nil {
			return err
		}

	}

	return nil
}

func buildVM(ctx *pulumi.Context, vm *VMConfig) error {
	if err := setResourceIDS(ctx, vm); err != nil {
		return err
	}
	var cloneArgs = &vsphere.VirtualMachineCloneArgs{}

	if vm.Template != "" {
		templateClone, err := vsphere.LookupVirtualMachine(ctx, &vsphere.LookupVirtualMachineArgs{
			Name:         vm.Template,
			DatacenterId: &vsphereIDs.dataCenterID,
		})

		if err != nil {
			return err
		}
		cloneArgs.LinkedClone = pulumi.Bool(false)
		cloneArgs.TemplateUuid = pulumi.String(templateClone.Uuid)
	}

	vmArgs := setVMOptions(vm)
	vmResult, err := vsphere.NewVirtualMachine(ctx, vm.Name, vmArgs)
	if err != nil {
		return err
	}

	ctx.Export(fmt.Sprintf("%s-IP", vm.Name), vmResult.NetworkInterfaces.ApplyT(func(v *vsphere.VirtualMachineNetworkInterface) (string, error) {
		if v == nil || v.DeviceAddress == nil {
			return "", nil
		}
		return *v.DeviceAddress, nil
	}).(pulumi.StringOutput))

	return nil
}

