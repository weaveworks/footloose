package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/cluster"
)

var showCmd = &cobra.Command{
	Use:     "show [HOSTNAME]",
	Aliases: []string{"status"},
	Short:   "Show all running machines or a single machine with a given hostname.",
	Long: `Provides information about machines created by footloose in JSON or Table format.
Optionally, provide show with a hostname to look for a specific machine. Exp: 'show node0'.`,
	RunE: show,
	Args: cobra.MaximumNArgs(1),
}

var showOptions struct {
	output string
	config string
}

func init() {
	showCmd.Flags().StringVarP(&showOptions.config, "config", "c", Footloose, "Cluster configuration file")
	showCmd.Flags().StringVarP(&showOptions.output, "output", "o", "table", "Output formatting options: {json,table}.")
	footloose.AddCommand(showCmd)
}

// show will show all machines in a given cluster.
func show(cmd *cobra.Command, args []string) error {
	c, err := cluster.NewFromFile(configFile(showOptions.config))
	if err != nil {
		return err
	}
	formatter, err := cluster.GetFormatter(showOptions.output)
	if err != nil {
		return err
	}
	machines, err := c.Inspect(args)
	if err != nil {
		return err
	}
	return formatter.Format(os.Stdout, machines)
}
