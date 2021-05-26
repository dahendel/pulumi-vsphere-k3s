package k3s

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

func Run(conf *VMConfig) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		if err := buildVM(ctx, conf); err != nil {
			return err
		}
		return nil
	}
}