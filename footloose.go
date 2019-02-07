package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var footloose = cobra.Command{
	Use:   "footloose",
	Short: "footloose - Container Machines",
}

func main() {
	log.SetFlags(0)

	if err := footloose.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
