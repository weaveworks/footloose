package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/weaveworks/footloose/pkg/config"

	"github.com/olekukonko/tablewriter"
)

// Formatter formats a slice of machines and outputs the result
// in a given format.
type Formatter interface {
	Format([]*Machine) error
	FormatSingle(Machine) error
}

// JSONFormatter formats a slice of machines into a JSON and
// outputs it to stdout.
type JSONFormatter struct{}

// TableFormatter formats a slice of machines into a colored
// table like output and prints that to stdout.
type TableFormatter struct{}

type port struct {
	Guest int `json:"guest"`
	Host  int `json:"host"`
}

const (
	// Stopped status of a machine
	Stopped = "Stopped"
	// Running status of a machine
	Running = "Running"
)

type status struct {
	Name     string          `json:"name"`
	State    string          `json:"state"`
	Spec     *config.Machine `json:"spec,omitempty"`
	Ports    []port          `json:"ports"`
	Hostname string          `json:"hostname"`
	Image    string          `json:"image"`
	Command  string          `json:"cmd"`
}

// Format will output to stdout in JSON format.
func (JSONFormatter) Format(machines []*Machine) error {
	statuses := make([]status, 0)
	for _, m := range machines {
		s := status{}
		s.Hostname = m.Hostname()
		s.Name = m.ContainerName()
		s.Image = m.spec.Image
		s.Command = m.spec.Cmd
		s.Spec = nil
		state := "Stopped"
		if m.IsRunning() {
			state = "Running"
		}
		s.State = state
		var ports []port
		for k, v := range m.ports {
			p := port{
				Host:  v,
				Guest: k,
			}
			ports = append(ports, p)
		}
		if len(ports) < 1 {
			for _, p := range m.spec.PortMappings {
				ports = append(ports, port{Host: int(p.ContainerPort), Guest: 0})
			}
		}
		s.Ports = ports
		statuses = append(statuses, s)
	}
	m := struct {
		Machines []status `json:"machines"`
	}{
		Machines: statuses,
	}
	ms, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(ms))
	return nil
}

// FormatSingle is a json formatter for a single machine.
func (JSONFormatter) FormatSingle(m Machine) error {
	s := status{}
	s.Hostname = m.Hostname()
	s.Name = m.ContainerName()
	s.Spec = m.spec
	s.Image = s.Spec.Image
	s.Command = s.Spec.Cmd
	state := Stopped
	if m.IsRunning() {
		state = Running
	}
	s.State = state
	var ports []port
	for k, v := range m.ports {
		p := port{
			Host:  v,
			Guest: k,
		}
		ports = append(ports, p)
	}
	s.Ports = ports
	ms, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", ms)
	return nil
}

// Format will output to stdout in table format.
func (TableFormatter) Format(machines []*Machine) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Hostname", "Ports", "Image", "Cmd", "State"})
	for _, m := range machines {
		state := Stopped
		if m.IsRunning() {
			state = Running
		}
		var ports []string
		for k, v := range m.ports {
			p := fmt.Sprintf("%d->%d", k, v)
			ports = append(ports, p)
		}
		if len(ports) < 1 {
			for _, p := range m.spec.PortMappings {
				port := fmt.Sprintf("%d->%d", p.ContainerPort, 0)
				ports = append(ports, port)
			}
		}
		ps := strings.Join(ports, ",")
		table.Append([]string{m.ContainerName(), m.Hostname(), ps, m.spec.Image, m.spec.Cmd, state})
	}
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.Render()
	return nil
}

// FormatSingle is a table formatter for a single machine.
func (TableFormatter) FormatSingle(machine Machine) error {
	jsonFormatter := JSONFormatter{}
	return jsonFormatter.FormatSingle(machine)
}

func getFormatter(output string) (Formatter, error) {
	var formatter Formatter
	switch output {
	case "json":
		formatter = new(JSONFormatter)
	case "table":
		formatter = new(TableFormatter)
	default:
		return nil, fmt.Errorf("unknown formatter '%s'", output)
	}
	return formatter, nil
}
