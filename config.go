package main

import (
	"github.com/spf13/cobra"

	"github.com/dlespiau/footloose/pkg/cluster"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a cluster configuration",
	RunE:  handleConfig,
}

var configOptions struct {
	file string
}

func init() {
	configCmd.Flags().StringVarP(&configOptions.file, "config", "c", Footloose, "Cluster configuration file")

	name := &defaultClusterSpec.Name
	configCmd.PersistentFlags().StringVarP(name, "name", "n", *name, "Name of the cluster")

	replicas := &defaultClusterSpec.Templates[0].Count
	configCmd.PersistentFlags().IntVarP(replicas, "replicas", "r", *replicas, "Number of machine replicas to config")

	footloose.AddCommand(configCmd)
}

func handleConfig(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(defaultClusterSpec)
	return cluster.Save(Footloose)
}
