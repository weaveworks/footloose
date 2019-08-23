package client

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/footloose/pkg/api"
	"github.com/weaveworks/footloose/pkg/config"
)

type env struct {
	server *httptest.Server
	client Client
}

func (e *env) Close() {
	e.server.Close()
}

func newEnv() *env {
	// Create an API server
	server := httptest.NewUnstartedServer(nil)
	baseURI := "http://" + server.Listener.Addr().String()
	api := api.New(baseURI)
	server.Config.Handler = api.Router()
	server.Start()

	return &env{
		server: server,
		client: Client{
			baseURI: server.URL,
			client:  server.Client(),
		},
	}
}

func TestCreateDeleteCluster(t *testing.T) {
	env := newEnv()
	defer env.Close()

	err := env.client.CreateCluster(&config.Cluster{
		Name:       "testcluster",
		PrivateKey: "testcluster-key",
	})
	assert.NoError(t, err)

	err = env.client.DeleteCluster("testcluster")
	assert.NoError(t, err)
}

func TestCreateDeleteMachine(t *testing.T) {
	env := newEnv()
	defer env.Close()

	err := env.client.CreateCluster(&config.Cluster{
		Name:       "testcluster",
		PrivateKey: "testcluster-key",
	})
	assert.NoError(t, err)

	err = env.client.CreateMachine("testcluster", &config.Machine{
		Name:  "testmachine",
		Image: "quay.io/footloose/centos7:latest",
		PortMappings: []config.PortMapping{
			{ContainerPort: 22},
		},
	})
	assert.NoError(t, err)

	err = env.client.DeleteMachine("testcluster", "testmachine")
	assert.NoError(t, err)

	err = env.client.DeleteCluster("testcluster")
	assert.NoError(t, err)
}
