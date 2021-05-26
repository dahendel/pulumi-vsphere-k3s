package k3s

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
	Size            int
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
