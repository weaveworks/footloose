package cluster

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/weaveworks/footloose/pkg/config"
)

// Formatter formats a slice of machines and outputs the result
// in a given format.
type Formatter interface {
	Format(io.Writer, []*Machine) error
	FormatSingle(io.Writer, *Machine) error
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
	// NotCreated status of a machine
	NotCreated = "Not created"
	// Stopped status of a machine
	Stopped = "Stopped"
	// Running status of a machine
	Running = "Running"
)

// MachineStatus is the runtime status of a Machine.
type MachineStatus struct {
	Container       string            `json:"container"`
	State           string            `json:"state"`
	Spec            *config.Machine   `json:"spec,omitempty"`
	Ports           []port            `json:"ports"`
	Hostname        string            `json:"hostname"`
	Image           string            `json:"image"`
	Command         string            `json:"cmd"`
	RuntimeNetworks []*RuntimeNetwork `json:"runtimeNetworks,omitempty"`
}

// Format will output to stdout in JSON format.
func (JSONFormatter) Format(w io.Writer, machines []*Machine) error {
	var statuses []MachineStatus
	for _, m := range machines {
		statuses = append(statuses, *m.Status())
	}

	m := struct {
		Machines []MachineStatus `json:"machines"`
	}{
		Machines: statuses,
	}
	ms, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	ms = append(ms, '\n')
	_, err = w.Write(ms)
	return err
}

// FormatSingle is a json formatter for a single machine.
func (js JSONFormatter) FormatSingle(w io.Writer, m *Machine) error {
	status, err := json.MarshalIndent(m.Status(), "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(status)
	return err
}

type tableMachine struct {
	Container string
	Hostname  string
	Ports     string
	IP        string
	Image     string
	Cmd       string
	State     string
	Backend   string
}

func writeColumns(w io.Writer, cols []string) {
	fmt.Fprintln(w, strings.Join(cols, "\t"))
}

// Format will output to stdout in table format.
func (TableFormatter) Format(w io.Writer, machines []*Machine) error {
	const padding = 3
	table := tabwriter.NewWriter(w, 0, 0, padding, ' ', 0)
	writeColumns(table, []string{"NAME", "HOSTNAME", "PORTS", "IP", "IMAGE", "CMD", "STATE", "BACKEND"})
	for _, m := range machines {
		state := NotCreated
		if m.IsCreated() {
			state = Stopped
			if m.IsStarted() {
				state = Running
			}
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
			Container: m.ContainerName(),
			Hostname:  m.Hostname(),
			Ports:     ps,
			IP:        m.ip,
			Image:     m.spec.Image,
			Cmd:       m.spec.Cmd,
			State:     state,
			Backend:   m.spec.Backend,
		}
		writeColumns(table, []string{tm.Container, tm.Hostname, tm.Ports, tm.IP, tm.Image, tm.Cmd, tm.State, tm.Backend})
	}
	table.Flush()
	return nil
}

// FormatSingle is a table formatter for a single machine.
func (TableFormatter) FormatSingle(w io.Writer, machine *Machine) error {
	jsonFormatter := JSONFormatter{}
	return jsonFormatter.FormatSingle(w, machine)
}

func GetFormatter(output string) (Formatter, error) {
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
