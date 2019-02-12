package main

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage cluster configuration",
}

func init() {
	footloose.AddCommand(configCmd)
}
