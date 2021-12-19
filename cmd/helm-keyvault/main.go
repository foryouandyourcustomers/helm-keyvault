package main

import (
	"errors"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func main() {

	fkv := cli.StringFlag{
		Name:     "keyvault",
		Aliases:  []string{"kv"},
		Usage:    "Name of the keyvault",
		Required: true,
	}

	fse := cli.StringFlag{
		Name:     "secret",
		Aliases:  []string{"s"},
		Usage:    "Name of the secret",
		Required: true,
	}

	fke := cli.StringFlag{
		Name:     "key",
		Aliases:  []string{"k"},
		Usage:    "Name of the key",
		Required: true,
	}

	fve := cli.StringFlag{
		Name:     "version",
		Aliases:  []string{"v"},
		Usage:    "Key or secret version",
		Required: false,
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "Decode keyvault secret or encrypted file, print result to stdout for usage with the helm downloader plugin",
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
				Name:    "secrets",
				Aliases: []string{"s", "secret"},
				Usage:   "get secret, put secrets, download secrets as values.yaml files",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "Get base64 encoded secret and decode it",
						Flags: []cli.Flag{
							&fkv,
							&fse,
							&fve,
						},
						Action: func(c *cli.Context) error {
							return cmd.GetSecret(c.String("keyvault"), c.String("secret"), c.String("version"))
						},
					},
					{
						Name:  "list",
						Usage: "List all secrets in the keyvault",
						Flags: []cli.Flag{
							&fkv,
						},
						Action: func(c *cli.Context) error {
							return cmd.ListSecrets(c.String("keyvault"))
						},
					},
					{
						Name:  "put",
						Usage: "Read file, base64 encode it and put it into keyvault",
						Flags: []cli.Flag{
							&fkv,
							&fse,
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Usage:    "path to file to encode and upload to keyvault",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return cmd.PutSecret(c.String("keyvault"), c.String("secret"), c.String("file"))
						},
					},
				},
			},
			{
				Name:    "keys",
				Aliases: []string{"k", "key"},
				Usage:   "create, export and list keys",
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create an azure keyvault key which can be used for local file encryption",
						Flags: []cli.Flag{
							&fkv,
							&fke,
						},
						Action: func(c *cli.Context) error {
							return cmd.CreateKey(c.String("keyvault"), c.String("key"))
						},
					},
					{
						Name:  "backup",
						Usage: "Backup azure keyvault key. The created backup can be imported into a keyvault and reused",
						Flags: []cli.Flag{
							&fkv,
							&fke,
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Usage:    "Backup filename",
								Required: false,
							},
						},
						Action: func(c *cli.Context) error {
							fn := c.String("file")
							if fn == "" {
								fn = fmt.Sprintf("%s.pem", strings.ToUpper(c.String("key")))
							}
							return cmd.BackupKey(c.String("keyvault"), c.String("key"), fn)
						},
					},
					{
						Name:  "list",
						Usage: "List all keys in the keyvault",
						Flags: []cli.Flag{
							&fkv,
						},
						Action: func(c *cli.Context) error {
							return cmd.ListKeys(c.String("keyvault"))
						},
					},
				},
			},
			{
				Name:    "files",
				Aliases: []string{"f", "file"},
				Usage:   "Encrypt and decrypt local files",
				Subcommands: []*cli.Command{
					{
						Name:  "encrypt",
						Usage: "Encrypt given file with given keyvault key",
						Flags: []cli.Flag{
							&fkv,
							&fke,
							&fve,
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Usage:    "BFile to encrypt",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return cmd.EncryptFile(c.String("keyvault"), c.String("key"), c.String("version"), c.String("file"))
						},
					},
					{
						Name:  "decrypt",
						Usage: "Decrypt the given file with the stored keyvault key",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "keyvault",
								Aliases:  []string{"kv"},
								Usage:    "Use alternate keyvault for decryption",
								Required: false,
								EnvVars:  []string{"KEYVAULT"},
							},
							&cli.StringFlag{
								Name:     "key",
								Aliases:  []string{"k"},
								Usage:    "Use alternate key for decryption",
								Required: false,
								EnvVars:  []string{"KEY"},
							},

							&cli.StringFlag{
								Name:     "version",
								Aliases:  []string{"v"},
								Usage:    "Use alternate version for decryption",
								Required: false,
								EnvVars:  []string{"VERSION"},
							},
							&cli.StringFlag{
								Name:     "file",
								Aliases:  []string{"f"},
								Usage:    "File to encrypt",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return cmd.DecryptFile(c.String("keyvault"), c.String("key"), "", c.String("file"))
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
