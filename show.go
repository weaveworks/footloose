package main

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/cluster"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show all running machines",
	RunE:  show,
	Args: cobra.MaximumNArgs(1),
}

var showOptions struct {
	output string
	config string
	all    bool
}

func init() {
	showCmd.Flags().StringVarP(&showOptions.config, "config", "c", Footloose, "Cluster configuration file")
	showCmd.Flags().StringVarP(&showOptions.output, "output", "o", "table", "Output options")
	showCmd.Flags().BoolVar(&showOptions.all, "all", false, "show all footloose created machines in every cluster.")
	footloose.AddCommand(showCmd)
}

// show will show all machines in a given cluster.
// if --all option is provided it will show every machine created by
// footloose no matter what cluster they are in.
func show(cmd *cobra.Command, args []string) error {
	c, err := cluster.NewFromFile(showOptions.config)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return c.Inspect(args[0])
	}
	return c.Show(showOptions.all, showOptions.output)
}
