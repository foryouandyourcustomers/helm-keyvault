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
	kv, err := structs.NewKeyvault(secret.Id.GetKeyvault())
	if err != nil {
		return "", err
	}

	secret.KeyVault = &kv
	secret.Name = secret.Id.GetName()
	secret.Version = secret.Id.GetVersion()

	secret, err = secret.Get()
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
	// parse uri to get file path
	parsed, err := url.Parse(u.uri)
	if err != nil {
		return "", err
	}

	var encfile structs.EncryptedFile
	encfile, err = encfile.LoadEncryptedFile(fmt.Sprintf("%s%s", parsed.Host, parsed.Path))
	if err != nil {
		return "", err
	}

	// decrypt the given data
	encfile.EncodedData, err = encfile.DecryptData(encfile.Kid.GetKeyvault(), encfile.Kid.GetName(), encfile.Kid.GetVersion())
	if err != nil {
		return "", err
	}

	value, err := encfile.GetDecodedString()
	if err != nil {
		return "", err
	}

	return value, nil
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
