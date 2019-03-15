package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print footloose version",
	Run:   showVersion,
}

func init() {
	footloose.AddCommand(versionCmd)
}

var version = "git"

func showVersion(cmd *cobra.Command, args []string) {
	fmt.Println("version:", version)
}
