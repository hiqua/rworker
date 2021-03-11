package main

import (
	"log"
	"os"

	"github.com/hiqua/rworker/internal/server"
	"github.com/urfave/cli/v2"
)

func serveFromContext(c *cli.Context) error {
	address := c.String("address")
	serverKey := c.String("serverKey")
	serverCertificate := c.String("serverCertificate")
	clientCertificates := c.StringSlice("clientCertificates")
	return server.Serve(address, serverCertificate, serverKey, clientCertificates...)
}

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:      "clientCertificates",
				Aliases:   []string{"cc"},
				Usage:     "Specify the list of client certificates to consider as users.",
				Required:  true,
				TakesFile: true,
			},
			&cli.PathFlag{
				Name:      "address",
				Aliases:   []string{"addr"},
				Usage:     "Specify the address to listen to (e.g. 'localhost:8443').",
				Required:  true,
				TakesFile: false,
			},
			&cli.PathFlag{
				Name:      "serverKey",
				Aliases:   []string{"sk"},
				Usage:     "Specify the server key.",
				Required:  true,
				TakesFile: true,
			},
			&cli.PathFlag{
				Name:      "serverCertificate",
				Aliases:   []string{"sc"},
				Usage:     "Specify the server certificate.",
				Required:  true,
				TakesFile: true,
			},
		},
		Action: serveFromContext,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
