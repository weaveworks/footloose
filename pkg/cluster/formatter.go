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

type status struct {
	Container       string            `json:"container"`
	State           string            `json:"state"`
	Spec            *config.Machine   `json:"spec,omitempty"`
	Ports           []port            `json:"ports"`
	Hostname        string            `json:"hostname"`
	Image           string            `json:"image"`
	Command         string            `json:"cmd"`
	RuntimeNetworks []*RuntimeNetwork `json:"runtime_networks,omitempty"`
}

// Format will output to stdout in JSON format.
func (JSONFormatter) Format(w io.Writer, machines []*Machine) error {
	var statuses []status
	for _, m := range machines {
		s := status{}
		s.Hostname = m.Hostname()
		s.Container = m.ContainerName()
		s.Image = m.spec.Image
		s.Command = m.spec.Cmd
		s.Spec = m.spec
		state := NotCreated
		if m.IsCreated() {
			state = Stopped
			if m.IsStarted() {
				state = Running
			}
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
		s.RuntimeNetworks = m.runtimeNetworks

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
	w.Write(ms)
	return nil
}

// FormatSingle is a json formatter for a single machine.
func (js JSONFormatter) FormatSingle(w io.Writer, m *Machine) error {
	return js.Format(w, []*Machine{m})
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
