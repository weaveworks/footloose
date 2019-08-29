package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/weaveworks/footloose/pkg/api"
	"github.com/weaveworks/footloose/pkg/cluster"
	"github.com/weaveworks/footloose/pkg/config"
)

// Client is a object able to talk a remote footloose API server.
type Client struct {
	baseURI *url.URL
	client  *http.Client
}

// New creates a new Client.
func New(baseURI string) *Client {
	u, err := url.Parse(baseURI)
	if err != nil {
		panic(err)
	}
	return &Client{
		baseURI: u,
		client:  &http.Client{},
	}
}

func (c *Client) uriFromPath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return c.baseURI.ResolveReference(u).String()
}

func (c *Client) publicKeyURI(name string) string {
	return c.uriFromPath(fmt.Sprintf("/api/keys/%s", name))
}

func (c *Client) clusterURI(name string) string {
	return c.uriFromPath(fmt.Sprintf("/api/clusters/%s", name))
}

func (c *Client) machineURI(clusterName, name string) string {
	return c.uriFromPath(fmt.Sprintf("/api/clusters/%s/machines/%s", clusterName, name))
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

func (c *Client) get(uri string, data interface{}) error {
	req, err := http.NewRequest("GET", uri, http.NoBody)
	if err != nil {
		return errors.Wrapf(err, "new GET request to %q", uri)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "http request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(apiError(resp), "GET status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return errors.Errorf("could not decode GET response: %v", err)
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

// CreatePublicKey creates a new public key.
func (c *Client) CreatePublicKey(def *config.PublicKey) error {
	return c.create(c.uriFromPath("/api/keys"), def)
}

// GetPublicKey retrieves a public key.
func (c *Client) GetPublicKey(name string) (*config.PublicKey, error) {
	data := config.PublicKey{}
	err := c.get(c.publicKeyURI(name), &data)
	return &data, err
}

// DeletePublicKey deletes a public key.
func (c *Client) DeletePublicKey(name string) error {
	return c.delete(c.publicKeyURI(name))
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

// GetMachine retrieves the machine details.
//
// XXX: This API isn't stable and will change in the future as we refine what
// the machine spec and status objects should be.
func (c *Client) GetMachine(clusterName, machine string) (*cluster.MachineStatus, error) {
	status := cluster.MachineStatus{}
	err := c.get(c.machineURI(clusterName, machine), &status)
	return &status, err
}

// DeleteMachine deletes a machine.
func (c *Client) DeleteMachine(cluster, machine string) error {
	return c.delete(c.machineURI(cluster, machine))
}
