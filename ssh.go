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

func init() {
	footloose.AddCommand(sshCmd)
}

func ssh(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(clusterSpec)
	return cluster.SSH(args[0])
}
