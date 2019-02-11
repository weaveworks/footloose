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

	name := &defaultConfig.Cluster.Name
	configCmd.PersistentFlags().StringVarP(name, "name", "n", *name, "Name of the cluster")

	private := &defaultConfig.Cluster.PrivateKey
	configCmd.PersistentFlags().StringVarP(private, "key", "k", *private, "Name of the private and public key files")

	replicas := &defaultConfig.Machines[0].Count
	configCmd.PersistentFlags().IntVarP(replicas, "replicas", "r", *replicas, "Number of machine replicas")

	privileged := &defaultConfig.Machines[0].Spec.Privileged
	configCmd.PersistentFlags().BoolVar(privileged, "privileged", *privileged, "Create privileged containers")

	footloose.AddCommand(configCmd)
}

func handleConfig(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(defaultConfig)
	return cluster.Save(configOptions.file)
}
