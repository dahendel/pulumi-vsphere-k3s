// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package vsphere

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The `VappContainer` data source can be used to discover the ID of a
// vApp container in vSphere. This is useful to fetch the ID of a vApp container
// that you want to use to create virtual machines in using the
// `VirtualMachine` resource.
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi-vsphere/sdk/v3/go/vsphere"
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		opt0 := "dc1"
// 		datacenter, err := vsphere.LookupDatacenter(ctx, &vsphere.LookupDatacenterArgs{
// 			Name: &opt0,
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = vsphere.LookupVappContainer(ctx, &vsphere.LookupVappContainerArgs{
// 			DatacenterId: datacenter.Id,
// 			Name:         "vapp-container-1",
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupVappContainer(ctx *pulumi.Context, args *LookupVappContainerArgs, opts ...pulumi.InvokeOption) (*LookupVappContainerResult, error) {
	var rv LookupVappContainerResult
	err := ctx.Invoke("vsphere:index/getVappContainer:getVappContainer", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getVappContainer.
type LookupVappContainerArgs struct {
	// The managed object reference
	// ID of the datacenter the vApp container is located in.
	DatacenterId string `pulumi:"datacenterId"`
	// The name of the vApp container. This can be a name or
	// path.
	Name string `pulumi:"name"`
}

// A collection of values returned by getVappContainer.
type LookupVappContainerResult struct {
	DatacenterId string `pulumi:"datacenterId"`
	// The provider-assigned unique ID for this managed resource.
	Id   string `pulumi:"id"`
	Name string `pulumi:"name"`
}
