package main

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
	NamePrefix       string
	Template         string
	NumCpus          int
	Cores            int
	Memory           int
	Networks         []*Network
	Disks            []Disk
	name             string
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

type AnsibleConfig struct {
	PrivateKey  string   `pulumi:"private_key"`
	// Playbooks is a list of paths to playbooks to execute
	Playbooks   []string `pulumi:"playbook"`
	// User for ansible to use upon connection
	AnsibleUser string   `pulumi:"ansible_user"`
}
