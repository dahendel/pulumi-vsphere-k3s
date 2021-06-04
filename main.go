package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		vsphereIDs = &vSphereIdData{}
		conf := &VMConfig{}
		pulumiConfig := config.New(ctx, "")
		pulumiConfig.RequireObject("vmConfig", conf)
		if err := buildVMs(ctx, conf); err != nil {
			return err
		}
		return nil
	})
}
