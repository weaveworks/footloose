package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into a machine",
	Args:  validateArgs,
	RunE:  ssh,
}

var sshOptions struct {
	config string
}

func init() {
	sshCmd.Flags().StringVarP(&sshOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(sshCmd)
}

func ssh(cmd *cobra.Command, args []string) error {
	cluster, err := cluster.NewFromFile(configFile(sshOptions.config))
	if err != nil {
		return err
	}
	var node string
	var username string
	if strings.Contains(args[0], "@") {
		items := strings.Split(args[0], "@")
		if len(items) != 2 {
			return fmt.Errorf("bad syntax for user@node: %v", items)
		}
		username = items[0]
		node = items[1]
	} else {
		node = args[0]
		machine, err := cluster.GetMachineByHostname(node)
		if err != nil {
			return fmt.Errorf("host not found: %s", node)
		}
		username = machine.User()
	}
	return cluster.SSH(node, username, args[1:]...)
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing machine name argument")
	}
	return nil
}
