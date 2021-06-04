package main

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

func newDiskArray(disk vsphere.GetVirtualMachineDisk) vsphere.VirtualMachineDiskArray {
	var diskArray vsphere.VirtualMachineDiskArray

	diskArray = append(diskArray, &vsphere.VirtualMachineDiskArgs{
		EagerlyScrub:    pulumi.BoolPtr(disk.EagerlyScrub),
		Label:           pulumi.StringPtr(disk.Label),
		Size:            pulumi.IntPtr(disk.Size),
		ThinProvisioned: pulumi.BoolPtr(disk.ThinProvisioned),
		UnitNumber:      pulumi.IntPtr(disk.UnitNumber),
	})
	//for i, d := range disks[1:] {
	//	diskArgs := &vsphere.VirtualMachineDiskArgs{
	//		DatastoreId:     pulumi.StringPtr(vsphereIDs.dataStoreID),
	//		EagerlyScrub:    pulumi.BoolPtr(d.EagerlyScrub),
	//		Label:           pulumi.StringPtr(fmt.Sprintf("%s-disk-%d.vmdk", namePrefix, i+1)),
	//		Size:            d.Size,
	//		ControllerType:  pulumi.StringPtr("scsi"),
	//		ThinProvisioned: pulumi.BoolPtr(d.ThinProvisioned),
	//		UnitNumber:      pulumi.IntPtr(i+1),
	//	}
	//	diskArray = append(diskArray, diskArgs)
	//}
	return diskArray
}

func newCloneArgs(templateID string) vsphere.VirtualMachineCloneArgs {
	return vsphere.VirtualMachineCloneArgs{
		Customize:    nil,
		LinkedClone:  pulumi.Bool(false),
		TemplateUuid: pulumi.String(templateID),
	}
}

func setVMOptions(conf *VMConfig, tmplId *vsphere.LookupVirtualMachineResult) *vsphere.VirtualMachineArgs {
	var diskArray vsphere.VirtualMachineDiskArray
	var networkArray vsphere.VirtualMachineNetworkInterfaceArray
	if conf.Disks != nil {
		diskArray = newDiskArray(tmplId.Disks[0])
	}

	if conf.Networks != nil {
		networkArray = newNetworkArray(conf.Networks)
	}

	clone := newCloneArgs(tmplId.Uuid)

	if conf.Folder == "" {
		conf.Folder = "/"
	}

	return &vsphere.VirtualMachineArgs{
		WaitForGuestIpTimeout: pulumi.Int(200),
		CpuHotAddEnabled:      pulumi.Bool(true),
		CpuHotRemoveEnabled:   pulumi.Bool(true),
		EnableDiskUuid:        pulumi.BoolPtr(conf.EnableDiskUuid),
		GuestId:               pulumi.StringPtr(tmplId.GuestId),
		Folder:                pulumi.StringPtr(conf.Folder),
		Memory:                pulumi.IntPtr(conf.Memory),
		MemoryHotAddEnabled:   pulumi.Bool(true),
		Name:                  pulumi.StringPtr(conf.name),
		NumCoresPerSocket:     pulumi.IntPtr(conf.Cores),
		NumCpus:               pulumi.IntPtr(conf.NumCpus),
		SyncTimeWithHost:      pulumi.Bool(true),
		ResourcePoolId:        pulumi.String(vsphereIDs.resourcePoolID),
		Disks:                 diskArray,
		NetworkInterfaces:     networkArray,
		Clone:                 clone,
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

func buildVMs(ctx *pulumi.Context, conf *VMConfig) error {
	if err := setResourceIDS(ctx, conf); err != nil {
		return err
	}

	var hosts []string
	for i := 1; i <= conf.Count; i++ {
		var vmName string
		vmName = fmt.Sprintf("%s%d", conf.NamePrefix, i+1)
		tmpl, err := vsphere.LookupVirtualMachine(ctx, &vsphere.LookupVirtualMachineArgs{
			Name:         conf.Template,
			DatacenterId: &vsphereIDs.dataCenterID,
		})

		if err != nil {
			return err
		}

		conf.name = vmName
		vmArgs := setVMOptions(conf, tmpl)
		vm, err := vsphere.NewVirtualMachine(ctx, vmName, vmArgs)
		if err != nil {
			return err
		}

		pulumi.Any(vm.DefaultIpAddress).ApplyT(func(args interface{}) string {
			hosts = append(hosts, args.(string))
			return args.(string)
		})
		ctx.Export(vmName, vm.DefaultIpAddress)

	}
	ctx.Export("ips", pulumi.ToStringArray(hosts))
	return nil
}
