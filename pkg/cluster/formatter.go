package cluster

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apcera/termtables"
	"github.com/weaveworks/footloose/pkg/config"
)

// Formatter formats a slice of machines and outputs the result
// in a given format.
type Formatter interface {
	Format([]*Machine) error
}

// JSONFormatter formats a slice of machines into a JSON and
// outputs it to stdout.
type JSONFormatter struct{}

// NormalFormatter formats a slice of machines into a colored
// table like output and prints that to stdout.
type NormalFormatter struct{}

type status struct {
	Name   string         `json:"name"`
	Ports  machineCache   `json:"ports"`
	Spec   config.Machine `json:"spec"`
	Status string         `json:"status"`
}

// Format will output to stdout in JSON format.
func (JSONFormatter) Format(machines []*Machine) error {
	var statuses []status
	for _, m := range machines {
		s := status{}
		s.Name = m.ContainerName()
		s.Ports = m.machineCache
		s.Spec = *m.spec
		state := "Stopped"
		if m.IsRunning() {
			state = "Running"
		}
		s.Status = state
		statuses = append(statuses, s)
	}
	m := struct {
		Machines []status `json:"machines"`
	}{
		Machines: statuses,
	}
	ms, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Printf("%s", ms)
	return nil
}

// Format will output to stdout in table format.
func (NormalFormatter) Format(machines []*Machine) error {
	table := termtables.CreateTable()
	table.AddHeaders("Name", "Ports", "State")
	for _, m := range machines {
		state := "Stopped"
		if m.IsRunning() {
			state = "Running"
		}
		table.AddRow(m.ContainerName(), "22", state)
	}
	fmt.Println(table.Render())
	return nil
}

func getFormatter(format string) (Formatter, error) {
	var formatter Formatter
	switch format {
	case "json":
		formatter = new(JSONFormatter)
	case "default":
		formatter = new(NormalFormatter)
	default:
		return nil, errors.New("unrecognised formatting method")
	}
	return formatter, nil
}
