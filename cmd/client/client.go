package main

import (
	"log"
	"os"

	"github.com/hiqua/rworker/internal/client"
	"github.com/urfave/cli/v2"
)

// TODO: could validate the job id given as argument (parses as valid uuid)
func main() {
	// TODO: the cli can always be nicer, .conf with default values, autocompletion, better error messages...
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      "address",
				Aliases:   []string{"addr"},
				Usage:     "Specify the server address (e.g. 'localhost:8443').",
				Required:  true,
				TakesFile: false,
			},
			&cli.PathFlag{
				Name:      "clientCert",
				Aliases:   []string{"cc"},
				Usage:     "Specify the folder containing the client's certificate.",
				Required:  true,
				TakesFile: true,
				Value:     "certs/client/cert.pem",
			},
			&cli.PathFlag{
				Name:      "clientKey",
				Aliases:   []string{"ck"},
				Usage:     "Specify the folder containing the client's key.",
				Required:  true,
				TakesFile: true,
				Value:     "certs/client/key.pem",
			},
			&cli.PathFlag{
				Name:      "serverCert",
				Aliases:   []string{"sc"},
				Usage:     "Specify the folder containing the server's certificate.",
				Required:  true,
				TakesFile: true,
				Value:     "certs/server/cert.pem",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add CMD [ARGS]...",
				Action:  rclient.AddJob,
			},
			{
				Name:    "status",
				Aliases: []string{"sta"},
				Usage:   "status JID",
				Action:  rclient.FetchStatus,
			},
			{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "log JID",
				Action:  rclient.FetchLog,
			},
			{
				Name:    "stop",
				Aliases: []string{"sto"},
				Usage:   "stop JID",
				Action:  rclient.StopJob,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
