package client

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/footloose/pkg/api"
	"github.com/weaveworks/footloose/pkg/cluster"
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
	keyStore := cluster.NewKeyStore(".")
	api := api.New(baseURI, keyStore)
	server.Config.Handler = api.Router()
	server.Start()

	u, _ := url.Parse(server.URL)

	return &env{
		server: server,
		client: Client{
			baseURI: u,
			client:  server.Client(),
		},
	}
}

const publicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDT3IG4sRIpLaoAtQSXBYaVLZTXh3Pl95ONm9oe9+nJ08qrUOFEJuKMTnqSgbC+R6v3T6fcgu1HgZtQyqB15rlA5U6rybKEa631+2Y+STBdCtBover2/c59QqfEyXWoPeq0EWRCt/ixVJdcTZqxNpZQUBoUQAIl1T/+lqEsefI4H/fFCeuqDyZfjWQXpoIh8fTpYleS6rmzvKTBhxg149LdmI96mo8Wzh2nSuXxxrk4ItvjUkNP/+s/I1xBZ6OKkO5a1Ngjuv4Yi0HM3SwZcIEP4P8QnFJtTUZjz7NyyPUthJy7QPIRMmimCg+yyRwkMhnbb6bNY6QIbQmrRw4rbGyd31eY/xXXLk6DLVGaoacVD5VuPjSEVjn9lzgaQoO1HJLYnAfgJB+3L/eKG5C8iE4gwnNbKMazLr2iVa6VdeACqyzTyx3uv/4TY2Q3Aqq+LPzOda0nbeaeIaq6xpA1iBsdNM/j88SOGJtYufUngVMql7nZGsxHt4oEw0OOGtshWcR27bKMJsuOkghnHJzs9o9uRBvBStZFLpEyA6TEIeNfTn6Mzdag/T+0NeisXUKSEvrMaxEVAnX7uvkMr5UNUeT/TDbVhAtFHm4YDFEnSupmMsAKiuiTA+XhBuY+FzsGTDGcVZRj6ERZl6u0A+Oo8p/h7TizP3ct7dXVD02dmfJGAQ== cluster@footloose.mail"

func TestCreateDeletePublicKey(t *testing.T) {
	env := newEnv()
	defer env.Close()

	err := env.client.CreatePublicKey(&config.PublicKey{
		Name: "testpublickey",
		Key:  publicKey,
	})
	assert.NoError(t, err)

	data, err := env.client.GetPublicKey("testpublickey")
	assert.Equal(t, "testpublickey", data.Name)
	assert.Equal(t, publicKey, data.Key)
	assert.NoError(t, err)

	err = env.client.DeletePublicKey("testpublickey")
	assert.NoError(t, err)
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

	status, err := env.client.GetMachine("testcluster", "testmachine")
	assert.NoError(t, err)
	assert.Equal(t, "testmachine", status.Spec.Name)

	err = env.client.DeleteMachine("testcluster", "testmachine")
	assert.NoError(t, err)

	err = env.client.DeleteCluster("testcluster")
	assert.NoError(t, err)
}
