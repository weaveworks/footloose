package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/weaveworks/footloose/pkg/api"
	"github.com/weaveworks/footloose/pkg/config"
)

// Client is a object able to talk a remote footloose API server.
type Client struct {
	baseURI string
	client  *http.Client
}

// New creates a new Client.
func New(baseURI string) *Client {
	return &Client{
		baseURI: baseURI,
		client:  &http.Client{},
	}
}

func (c *Client) uriFromPath(path string) string {
	return c.baseURI + path
}

func (c *Client) clusterURI(name string) string {
	return fmt.Sprintf("%s/api/clusters/%s", c.baseURI, name)
}

func (c *Client) machineURI(clusterName, name string) string {
	return fmt.Sprintf("%s/api/clusters/%s/machines/%s", c.baseURI, clusterName, name)
}

func apiError(resp *http.Response) error {
	e := api.ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return errors.New("could not decode error response")
	}
	return errors.New(e.Error)
}

func (c *Client) create(uri string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "new POST request to %q", uri)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.Wrapf(apiError(resp), "POST status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) delete(uri string) error {
	req, err := http.NewRequest("DELETE", uri, http.NoBody)
	if err != nil {
		return errors.Wrapf(err, "new DELETE request to %q", uri)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "http request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(apiError(resp), "DELETE status %d", resp.StatusCode)
	}
	return nil
}

// CreateCluster creates a new cluster.
func (c *Client) CreateCluster(def *config.Cluster) error {
	return c.create(c.uriFromPath("/api/clusters"), def)
}

// DeleteCluster deletes a cluster and all its associated machines.
func (c *Client) DeleteCluster(name string) error {
	return c.delete(c.clusterURI(name))
}

// CreateMachine creates a new machine.
func (c *Client) CreateMachine(cluster string, def *config.Machine) error {
	return c.create(c.uriFromPath(fmt.Sprintf("/api/clusters/%s/machines", cluster)), def)
}

// DeleteMachine deletes a machine.
func (c *Client) DeleteMachine(cluster, machine string) error {
	return c.delete(c.machineURI(cluster, machine))
}
