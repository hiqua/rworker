package rclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/hiqua/rworker/internal/directory"
	"github.com/urfave/cli/v2"
)

// AddJob adds a job.
func AddJob(c *cli.Context) error {
	client, err := InitializeClientFromContext(c)
	if err != nil {
		return err
	}
	jobRequest := directory.JobRequest{Command: c.Args().First(), Arguments: c.Args().Tail()}
	bs, err := json.Marshal(jobRequest)
	if err != nil {
		return err
	}

	r, err := client.Post(getAddressFromContext(c)+"/job", "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}

	return PrintResponse(r.Body)
}

// FetchLog fetches the logs.
func FetchLog(c *cli.Context) error {
	// TODO: we could support fetching logs of several jobs at the same time
	client, err := InitializeClientFromContext(c)
	if err != nil {
		return err
	}
	r, err := client.Get(getAddressFromContext(c) + "/log/" + c.Args().First())
	if err != nil {
		return err
	}

	return PrintResponse(r.Body)
}

// FetchStatus fetches the status.
func FetchStatus(c *cli.Context) error {
	client, err := InitializeClientFromContext(c)
	if err != nil {
		return err
	}
	r, err := client.Get(getAddressFromContext(c) + "/job/" + c.Args().First())
	if err != nil {
		return err
	}

	return PrintResponse(r.Body)
}

// StopJob stops a job.
func StopJob(c *cli.Context) error {
	client, err := InitializeClientFromContext(c)
	if err != nil {
		return err
	}
	r, err := ClientDelete(client, getAddressFromContext(c)+"/stop/"+c.Args().First())
	if err != nil {
		return err
	}

	return PrintResponse(r.Body)
}

// InitializeClientFromContext sets up a client and the TLS certificates given from the CLI context.
func InitializeClientFromContext(c *cli.Context) (*http.Client, error) {
	clientCert, serverCert := c.String("clientCert"), c.String("serverCert")
	clientKey := c.String("clientKey")
	return InitializeClient(clientCert, clientKey, serverCert)
}

func getAddressFromContext(c *cli.Context) string {
	return "https://" + c.String("address")
}

// InitializeClient sets up a client and the TLS certificates.
func InitializeClient(clientCert, clientKey, serverCert string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(serverCert)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	return client, nil

}

// PrintResponse prints the response received from the server.
func PrintResponse(body io.ReadCloser) error {
	bs, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return nil
}

// ClientDelete sends a DELETE HTTP request
// TODO: is there no .Delete in the library?
func ClientDelete(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
