package main

import (
	"github.com/spf13/cobra"

	"github.com/dlespiau/footloose/pkg/cluster"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cluster",
	RunE:  create,
}

func init() {
	replicas := &clusterSpec.Templates[0].Count
	createCmd.PersistentFlags().IntVarP(replicas, "--replicas", "r", *replicas, "Number of machine replicas to create")
	footloose.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(clusterSpec)
	return cluster.Create()
}
