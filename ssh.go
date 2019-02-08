package main

import (
	"github.com/spf13/cobra"

	"github.com/dlespiau/footloose/pkg/cluster"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into a machine",
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
	cluster, err := cluster.NewFromFile(sshOptions.config)
	if err != nil {
		return err
	}
	return cluster.SSH(args[0])
}
