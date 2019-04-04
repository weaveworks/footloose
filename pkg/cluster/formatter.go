package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
	sortedNames := getSortedMachineNames(machines)
	statusMap := make(map[string]status, 0)
	for _, m := range machines {
		s := status{}
		s.Hostname = m.Hostname()
		s.Name = strings.TrimPrefix(m.ContainerName(), "/")
		s.Image = m.spec.Image
		s.Command = m.spec.Cmd
		s.Spec = m.spec
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
		statusMap[s.Name] = s
	}
	var statuses []status
	for _, name := range sortedNames {
		statuses = append(statuses, statusMap[name])
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

// getSortedMachineNames retrieves a sorted list of machine names.
func getSortedMachineNames(machines []*Machine) (names []string) {
	for _, m := range machines {
		names = append(names, strings.TrimPrefix(m.name, "/"))
	}
	sort.Strings(names)
	return
}

// FormatSingle is a json formatter for a single machine.
func (js JSONFormatter) FormatSingle(m Machine) error {
	return js.Format([]*Machine{&m})
}

type tableMachine struct {
	Name     string
	Hostname string
	Ports    string
	IP       string
	Image    string
	Cmd      string
	State    string
}

// Format will output to stdout in table format.
func (TableFormatter) Format(machines []*Machine) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Hostname", "Ports", "IP", "Image", "Cmd", "State"})
	machineMap := make(map[string]tableMachine)
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
				port := fmt.Sprintf("%d->%d", p.HostPort, p.ContainerPort)
				ports = append(ports, port)
			}
		}
		ps := strings.Join(ports, ",")
		tm := tableMachine{
			Name:     strings.TrimPrefix(m.ContainerName(), "/"),
			Hostname: m.Hostname(),
			Ports:    ps,
			IP:       m.ip,
			Image:    m.spec.Image,
			Cmd:      m.spec.Cmd,
			State:    state,
		}
		machineMap[tm.Name] = tm
	}
	sortedNames := getSortedMachineNames(machines)
	for _, name := range sortedNames {
		m := machineMap[name]
		table.Append([]string{m.Name, m.Hostname, m.Ports, m.IP, m.Image, m.Cmd, m.State})
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
