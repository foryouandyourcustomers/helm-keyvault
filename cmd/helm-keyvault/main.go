package main

import (
	"errors"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	sf := []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Retrieve secret with id https://<keyvault-name>.vault.azure.net/secrets/<secret>/<version>",
			Required: true,
		},
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Downlaod and decode secret from keyvault and print contents to stdout - use with helm downloader plugin",
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 4 {
						return errors.New("Please specify four arguments - command certFile keyFile caFile full-URL")
					}
					u := c.Args().Get(3)
					if len(u) <= 0 {
						return errors.New("full-URL argument missing")
					}
					return cmd.Download(u)

				},
			},
			{
				Name:    "secret",
				Aliases: []string{"s"},
				Usage:   "get secret, put secrets, download secrets as values.yaml files",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "Get base64 encoded secret and decode it",
						Flags: sf,
						Action: func(c *cli.Context) error {
							return cmd.GetSecret(c.String("id"))
						},
					},
					{
						Name:  "put",
						Usage: "Read file, base64 encode it and put it into keyvault",
						Flags: append(sf, &cli.StringFlag{
							Name:     "file",
							Usage:    "path to file to encode and upload to keyvault",
							Required: true,
						}),
						Action: func(c *cli.Context) error {
							return cmd.PutSecret(c.String("id"), c.String("file"))
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
