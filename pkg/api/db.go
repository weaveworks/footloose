package api

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/weaveworks/footloose/pkg/cluster"
)

type entry struct {
	cluster  *cluster.Cluster
	machines map[string]*cluster.Machine
}

type db struct {
	sync.Mutex

	clusters map[string]entry
}

func (db *db) init() {
	db.clusters = make(map[string]entry)
}

func (db *db) entry(name string) *entry {
	db.Lock()
	defer db.Unlock()

	entry, ok := db.clusters[name]
	if !ok {
		return nil
	}
	return &entry
}

func (db *db) cluster(name string) (*cluster.Cluster, error) {
	entry := db.entry(name)
	if entry == nil {
		return nil, errors.Errorf("unknown cluster '%s'", name)
	}
	return entry.cluster, nil
}

func (db *db) addCluster(name string, c *cluster.Cluster) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.clusters[name]; ok {
		return errors.Errorf("cluster '%s' has already been added", name)
	}
	db.clusters[name] = entry{
		cluster:  c,
		machines: make(map[string]*cluster.Machine),
	}
	return nil
}

func (db *db) removeCluster(name string) (*cluster.Cluster, error) {
	db.Lock()
	defer db.Unlock()

	var entry entry
	var ok bool
	if entry, ok = db.clusters[name]; !ok {
		return nil, errors.Errorf("unknown cluster '%s'", name)
	}
	// It is an error to remove the cluster from the db before removing all of its
	// machines.
	if len(entry.machines) != 0 {
		return nil, errors.Errorf("cluster has machines associated with it")
	}
	delete(db.clusters, name)
	return entry.cluster, nil
}

func (db *db) machine(clusterName, machineName string) (*cluster.Machine, error) {
	entry := db.entry(clusterName)
	if entry == nil {
		return nil, errors.Errorf("unknown cluster '%s'", clusterName)
	}

	db.Lock()
	defer db.Unlock()

	var m *cluster.Machine
	var ok bool
	if m, ok = entry.machines[machineName]; !ok {
		return nil, errors.Errorf("unknown machine '%s' for cluster '%s'", machineName, clusterName)
	}
	return m, nil
}

func (db *db) machines(clusterName string) ([]*cluster.Machine, error) {
	entry := db.entry(clusterName)
	if entry == nil {
		return nil, errors.Errorf("unknown cluster '%s'", clusterName)
	}

	db.Lock()
	defer db.Unlock()

	var machines []*cluster.Machine
	for _, m := range entry.machines {
		machines = append(machines, m)
	}
	return machines, nil
}

func (db *db) addMachine(cluster string, m *cluster.Machine) error {
	entry := db.entry(cluster)
	if entry == nil {
		return errors.Errorf("unknown cluster '%s'", cluster)
	}

	db.Lock()
	defer db.Unlock()

	// Hostname is really the machine unique name as we don't allow setting a
	// different hostname.
	if _, ok := entry.machines[m.Hostname()]; ok {
		return errors.Errorf("machine '%s' has already been added", m.Hostname())

	}
	entry.machines[m.Hostname()] = m
	return nil
}

func (db *db) removeMachine(clusterName, machineName string) (*cluster.Machine, error) {
	entry := db.entry(clusterName)
	if entry == nil {
		return nil, errors.Errorf("unknown cluster '%s'", clusterName)
	}

	db.Lock()
	defer db.Unlock()

	var m *cluster.Machine
	var ok bool
	if m, ok = entry.machines[machineName]; !ok {
		return nil, errors.Errorf("unknown machine '%s' for cluster '%s'", machineName, clusterName)
	}
	delete(entry.machines, machineName)
	return m, nil
}
