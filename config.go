package k3s

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type VMConfig struct {
	Datacenter       *string
	Cluster          string
	ClusterDatastore string
	DataStore        string
	ResourcePool     string
	Folder           string
	EnableDiskUuid   bool
	Count            int
	Name             string
	Template         string
	NumCpus          int
	Cores            int
	Memory           int
	Networks         []*Network
	Disks            []Disk
}

type Disk struct {
	Size            pulumi.Int
	Attach          bool
	EagerlyScrub    bool
	ThinProvisioned bool
}

type Network struct {
	AdapterType string
	NetworkName string
	MacAddress  string
	StaticMac   bool
	networkID   string
}
