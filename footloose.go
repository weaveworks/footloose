package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Footloose is the default name of the footloose file.
const Footloose = "footloose.yaml"

var footloose = &cobra.Command{
	Use:           "footloose",
	Short:         "footloose - Container Machines",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func main() {
	if err := footloose.Execute(); err != nil {
		log.Fatal(err)
	}
}
