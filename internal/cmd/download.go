package cmd

import (
	"encoding/base64"
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
	// parse uri to get file path
	parsed, err := url.Parse(u.uri)
	if err != nil {
		return "", err
	}

	encfile := structs.EncryptedFile{}
	err = encfile.LoadEncryptedFile(fmt.Sprintf("%s%s", parsed.Host, parsed.Path))
	if err != nil {
		return "", err
	}

	// retrieve keyvault information from loaded file
	kv, err := encfile.Kid.GetKeyvault()
	key, err := encfile.Kid.GetName()
	version, err := encfile.Kid.GetVersion()

	// decrypt the given data
	err = encfile.DecryptData(kv, key, version)
	if err != nil {
		return "", err
	}

	// parse the encoded chunks and return them as
	// a single string
	var value string
	for _, chunk := range encfile.EncodedData {
		c, err := base64.RawURLEncoding.DecodeString(chunk)
		if err != nil {
			return "", err
		}
		value = fmt.Sprintf("%s%s", value, string(c))
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
