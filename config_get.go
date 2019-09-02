package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/config"
)

var getConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config file information",
	RunE:  getConfig,
}

var getOptions struct {
	config string
}

func init() {
	getConfigCmd.Flags().StringVarP(&getOptions.config, "config", "c", Footloose, "Cluster configuration file")
	configCmd.AddCommand(getConfigCmd)
}

func getConfig(cmd *cobra.Command, args []string) error {
	c, err := config.NewConfigFromFile(configFile(getOptions.config))
	if err != nil {
		return err
	}
	var detail interface{}
	if len(args) > 0 {
		detail, err = config.GetValueFromConfig(args[0], c)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("Failed to get config detail")
		}
	} else {
		detail = c
	}
	if reflect.ValueOf(detail).Kind() != reflect.String {
		res, err := json.MarshalIndent(detail, "", "  ")
		if err != nil {
			log.Println(err)
			return fmt.Errorf("Cannot convert result to json")
		}
		fmt.Printf("%s", res)
	} else {
		fmt.Printf("%s", detail)
	}
	return nil
}
