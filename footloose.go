package main

import (
	"os"

	"github.com/spf13/cobra"
)

// Footloose is the default name of the footloose file.
const Footloose = "Footloose"

var footloose = &cobra.Command{
	Use:   "footloose",
	Short: "footloose - Container Machines",
}

func main() {
	if err := footloose.Execute(); err != nil {
		os.Exit(1)
	}
}
