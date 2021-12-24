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

func run(args []string) error {
	// flags used for cli commands
	flagKeyVault := cli.StringFlag{
		Name:     "keyvault",
		Aliases:  []string{"kv"},
		Usage:    "Name of the keyvault",
		Required: true,
		EnvVars:  []string{"KEYVAULT"},
	}

	flagSecret := cli.StringFlag{
		Name:     "secret",
		Aliases:  []string{"s"},
		Usage:    "Name of the secret",
		Required: true,
		EnvVars:  []string{"SECRET"},
	}

	flagKey := cli.StringFlag{
		Name:     "key",
		Aliases:  []string{"k"},
		Usage:    "Name of the key",
		Required: true,
		EnvVars:  []string{"KEY"},
	}

	flagVersion := cli.StringFlag{
		Name:     "version",
		Aliases:  []string{"v"},
		Usage:    "Key or secret version",
		Required: false,
		EnvVars:  []string{"VERSION"},
	}

	flagSecretFile := cli.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "path to file to encode and upload to keyvault as a secret",
		Required: true,
	}

	flagBackupFile := cli.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "Backup filename - defaults to \"[KEY|SECRET].pem\"",
		Required: false,
	}

	flagEncryptFile := cli.StringFlag{
		Name:     "file",
		Aliases:  []string{"f"},
		Usage:    "File to encrypt or decrypt with azure keyvault key",
		Required: true,
	}

	// the file decrypt option allows overwriting of the given keyvault, key and version
	// to do this we can specify optional values for keyvault, key and versio
	flagKeyVaultOptional := flagKeyVault
	flagKeyVaultOptional.Required = false
	flagKeyVaultOptional.Usage = "Use alternate keyvault for decryption"
	flagKeyOptional := flagKey
	flagKeyOptional.Required = false
	flagKeyOptional.Usage = "Use alternate key for decryption"
	flagVersionOptional := flagVersion
	flagVersionOptional.Usage = "Use alternate version for decryption"

	app := &cli.App{
		Name:  "helm-keyvault",
		Usage: "Manage Azure Keyvault secrets and keys to safely store and download helm charts.",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Sebastian Hutter",
				Email: "seh@foryouandyourcustomers.com",
			},
		},
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
							&flagKeyVault,
							&flagSecret,
							&flagVersion,
						},
						Action: func(c *cli.Context) error {
							return cmd.GetSecret(c.String("keyvault"), c.String("secret"), c.String("version"))
						},
					},
					{
						Name:  "list",
						Usage: "List all secrets in the keyvault",
						Flags: []cli.Flag{
							&flagKeyVault,
						},
						Action: func(c *cli.Context) error {
							return cmd.ListSecrets(c.String("keyvault"))
						},
					},
					{
						Name:  "put",
						Usage: "Read file, base64 encode it and put it into keyvault",
						Flags: []cli.Flag{
							&flagKeyVault,
							&flagSecret,
							&flagSecretFile,
						},
						Action: func(c *cli.Context) error {
							return cmd.PutSecret(c.String("keyvault"), c.String("secret"), c.String("file"))
						},
					},
					{
						Name:  "backup",
						Usage: "Backup azure keyvault secret. The created backup can be imported into a keyvault and reused",
						Flags: []cli.Flag{
							&flagKeyVault,
							&flagSecret,
							&flagBackupFile,
						},
						Action: func(c *cli.Context) error {
							fn := c.String("file")
							if fn == "" {
								fn = fmt.Sprintf("%s.pem", strings.ToUpper(c.String("secret")))
							}
							return cmd.BackupSecret(c.String("keyvault"), c.String("secret"), fn)
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
							&flagKeyVault,
							&flagKey,
						},
						Action: func(c *cli.Context) error {
							return cmd.CreateKey(c.String("keyvault"), c.String("key"))
						},
					},
					{
						Name:  "backup",
						Usage: "Backup azure keyvault key. The created backup can be imported into a keyvault and reused",
						Flags: []cli.Flag{
							&flagKeyVault,
							&flagKey,
							&flagBackupFile,
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
							&flagKeyVault,
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
							&flagKeyVault,
							&flagKey,
							&flagVersion,
							&flagEncryptFile,
						},
						Action: func(c *cli.Context) error {
							return cmd.EncryptFile(c.String("keyvault"), c.String("key"), c.String("version"), c.String("file"))
						},
					},
					{
						Name:  "decrypt",
						Usage: "Decrypt the given file with the stored keyvault key",
						Flags: []cli.Flag{
							&flagKeyVaultOptional,
							&flagKeyOptional,
							&flagVersionOptional,
							&flagEncryptFile,
						},
						Action: func(c *cli.Context) error {
							return cmd.DecryptFile(c.String("keyvault"), c.String("key"), "", c.String("file"))
						},
					},
				},
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
