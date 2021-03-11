package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"testing"

	rclient "github.com/hiqua/rworker/internal/client"
	"github.com/hiqua/rworker/internal/directory"
	"github.com/hiqua/rworker/internal/server"
	"github.com/stretchr/testify/assert"
)

// TODO: possibly flaky tests if server startup is slow
func init() {
	go func() {
		log.Fatal(server.Serve("127.0.0.1:58443",
			"certs/server/cert.pem",
			"certs/server/key.pem",
			"certs/client/cert.pem",
		))
	}()
}

func TestCanAddJobAsKnownClient(t *testing.T) {
	cc := filepath.Join("certs", "client", "cert.pem")
	ck := filepath.Join("certs", "client", "key.pem")
	sc := filepath.Join("certs", "server", "cert.pem")
	client, err := rclient.InitializeClient(cc, ck, sc)
	if assert.Nil(t, err) {
		err = submitNewJobRequest(client)
		assert.Nil(t, err)
	}
}

func TestCannotAddJobAsUnknownClient(t *testing.T) {
	cc := filepath.Join("certs", "badclient", "cert.pem")
	ck := filepath.Join("certs", "badclient", "key.pem")
	sc := filepath.Join("certs", "server", "cert.pem")
	client, err := rclient.InitializeClient(cc, ck, sc)
	if assert.Nil(t, err) {
		err = submitNewJobRequest(client)
		assert.NotNil(t, err)
	}
}

func submitNewJobRequest(client *http.Client) error {
	jobRequest := directory.JobRequest{Command: "ls", Arguments: []string{"/"}}
	bs, err := json.Marshal(jobRequest)
	if err != nil {
		return err
	}

	_, err = client.Post("https://localhost:58443/job", "application/json", bytes.NewReader(bs))
	return err
}
