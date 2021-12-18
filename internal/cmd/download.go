package cmd

import (
	"errors"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"net/url"
	"strings"
)

// interface for different uri download functions
type generalUri interface {
	download() (string, error)
}

// keyvaultUri - represents an keyvault+secret uri
type keyvaultUri struct {
	uri string
}

func (u *keyvaultUri) download() (string, error) {

	// replace the scheme - keyvault+secret(s):// with https
	uri := strings.Replace(u.uri, "keyvault+secrets://", "https://", 1)
	uri = strings.Replace(uri, "keyvault+secret://", "https://", 1)

	// create secrets struct
	secret := structs.Secret{Id: structs.KeyvaultObjectId(uri)}

	kv, _ := secret.Id.GetKeyvault()
	n, _ := secret.Id.GetName()
	v, _ := secret.Id.GetVersion()
	secret.KeyVault = kv
	secret.Name = n
	secret.Version = v

	println(secret.Id)
	println(kv)
	println(n)
	println(v)

	_, err := secret.Get()
	if err != nil {
		return "", err
	}

	val, err := secret.Decode()
	if err != nil {
		return "", err
	}

	return val, nil
}

// fileUri - represents an keyvault+file uri
type fileUri struct {
	uri string
}

func (u *fileUri) download() (string, error) {
	return "abc", nil
}

// DownloadSecret - Download and decode secret to be used as downloader plugin
func Download(uri string) error {

	u, err := parseUri(uri)
	if err != nil {
		return err
	}

	result, err := u.download()
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func parseUri(uri string) (generalUri, error) {

	u, err := url.Parse(uri[len("keyvault+"):])
	if err != nil {
		return nil, err
	}

	if (u.Scheme == "secret") || (u.Scheme == "secrets") {
		return &keyvaultUri{uri}, nil
	}
	if (u.Scheme == "file") || (u.Scheme == "files") {
		return &fileUri{uri}, nil
	}

	return nil, errors.New("Unknown download uri received")
}
