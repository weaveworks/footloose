package main

import (
	"os"

	"github.com/spf13/cobra"
)

var footloose = cobra.Command{
	Use:   "footloose",
	Short: "footloose - Container Machines",
}

func main() {
	if err := footloose.Execute(); err != nil {
		os.Exit(1)
	}
}
