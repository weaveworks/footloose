package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cluster configuration",
	RunE:  configCreate,
}

var configCreateOptions struct {
	file string
}

func init() {
	configCreateCmd.Flags().StringVarP(&configCreateOptions.file, "config", "c", Footloose, "Cluster configuration file")

	name := &defaultConfig.Cluster.Name
	configCreateCmd.PersistentFlags().StringVarP(name, "name", "n", *name, "Name of the cluster")

	private := &defaultConfig.Cluster.PrivateKey
	configCreateCmd.PersistentFlags().StringVarP(private, "key", "k", *private, "Name of the private and public key files")

	replicas := &defaultConfig.Machines[0].Count
	configCreateCmd.PersistentFlags().IntVarP(replicas, "replicas", "r", *replicas, "Number of machine replicas")

	image := &defaultConfig.Machines[0].Spec.Image
	configCreateCmd.PersistentFlags().StringVarP(image, "image", "i", *image, "Docker image to use in the containers")

	privileged := &defaultConfig.Machines[0].Spec.Privileged
	configCreateCmd.PersistentFlags().BoolVar(privileged, "privileged", *privileged, "Create privileged containers")

	configCmd.AddCommand(configCreateCmd)
}

// configExists checks whether a configuration file has already been created.
// Returns false if not true if it already exists.
func configExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || os.IsPermission(err) {
		return false
	}
	return !info.IsDir()
}

func configCreate(cmd *cobra.Command, args []string) error {
	cluster := cluster.New(defaultConfig)
	if configExists(configCreateOptions.file) {
		message := fmt.Sprintf("Configuration file at %s already exists...", configCreateOptions.file)
		return errors.New(message)
	}
	return cluster.Save(configCreateOptions.file)
}
