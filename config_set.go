package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/weaveworks/footloose/pkg/cluster"

	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/config"
)

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set value to configuration",
	RunE:  configSet,
}

var setOptions struct {
	config string
}

var setDryRun bool

func init() {
	configSetCmd.Flags().StringVarP(&setOptions.config, "config", "c", Footloose, "Cluster configuration file")
	configSetCmd.Flags().BoolVar(&setDryRun, "dry-run", defaultDryRun, "Dry run changes")
	configCmd.AddCommand(configSetCmd)
}

func configSet(cmd *cobra.Command, args []string) error {
	conf, cstr := getConfigAndCluster()
	checkSetRequirements(args, cstr, conf)
	err := config.SetValueToConfig(args[0], conf, config.ClarifyArg(args[1]))
	if err != nil {
		log.Fatalln(err)
	}
	handleSetResponse(conf)
	return nil
}

func getConfigAndCluster() (*config.Config, *cluster.Cluster) {
	conf, err := config.NewConfigFromFile(setOptions.config)
	if err != nil {
		log.Fatalln(err)
	}
	cstr, err := cluster.New(*conf)
	if err != nil {
		log.Fatalln(err)
	}
	return conf, cstr
}

func checkSetRequirements(args []string, c *cluster.Cluster, conf *config.Config) {
	if len(args) != 2 {
		log.Fatalln("set command needs 2 args")
	}
	err := config.IsSetValueValid(args[0], args[1])
	if err != nil {
		log.Fatalln(err)
	}
	num, err := c.CountCreatedMachine()
	if err != nil {
		log.Fatalln(err)
	}
	if num > 0 {
		log.Fatalln("cannot change config, please delete your machines before any change")
	}
}

func handleSetResponse(conf *config.Config) {
	if setDryRun == true {
		res, err := json.MarshalIndent(conf, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s", res)
	} else {
		conf.Save(setOptions.config)
	}
}
