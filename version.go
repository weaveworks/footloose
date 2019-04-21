package main

import (
	"fmt"
	"strings"

	"github.com/weaveworks/footloose/pkg/cluster"
	release "github.com/weaveworks/footloose/pkg/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print footloose version",
	Run:   showVersion,
}

func init() {
	versionCmd.Flags().BoolVarP(&verbosity, "verbosity", "v", false, "Verbosity commandline calls")
	footloose.AddCommand(versionCmd)
}

var version = "git"

func showVersion(cmd *cobra.Command, args []string) {
	cluster.GetCommanderInstance().SetVerbosity(verbosity)
	fmt.Println("version:", version)
	release, err := release.FindLastRelease()
	if err != nil {
		fmt.Println("Failed to check for new versions")
	}
	if strings.Compare(version, *release.TagName) != 0 {
		fmt.Printf("New version %v is available. More informations on: %v\n", *release.TagName, *release.HTMLURL)
	}
}
