package main

import (
	"context"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"os"
	"strings"
)

func runAnsiblePlaybook(hosts []string, ansibleConfig *AnsibleConfig) {
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "ssh",
		User:       ansibleConfig.AnsibleUser,
		PrivateKey: ansibleConfig.PrivateKey,
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: strings.Join(hosts, ","),
	}


	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks:         ansibleConfig.Playbooks,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
	}

	os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "false")

	err := playbook.Run(context.TODO())
	if err != nil {
		panic(err)
	}
}
