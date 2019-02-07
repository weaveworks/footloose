package main

import (
	"github.com/spf13/cobra"

	"github.com/dlespiau/footloose/pkg/cluster"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster",
	RunE:  delete,
}

func init() {
	footloose.AddCommand(deleteCmd)
}

func delete(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(clusterSpec)
	return cluster.Delete()
}
