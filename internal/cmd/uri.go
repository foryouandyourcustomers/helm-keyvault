package cmd

import (
	"encoding/json"
	"errors"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"net/url"
	"strings"
)

type generalUri interface {
	download() (string, error)
}

// secretUri - represents an keyvault+secret uri
type secretUri struct {
	Keyvault string
	Name     string
	Version  string
}

func (u *secretUri) download() (string, error) {

	//get secret from keyvault and print it
	secret, err := keyvault.GetSecret(u.Keyvault, u.Name, u.Version)
	if err != nil {
		return "", err
	}
	// decode secret and print
	dec, err := secret.Decode()
	if err != nil {
		return "", err
	}
	return dec, nil
}

// upload - upload base64 encoded content to keyvault
func (u *secretUri) upload(c string) (string, error) {

	s, err := keyvault.PutSecret(u.Keyvault, u.Name, c)
	if err != nil {
		return "", err
	}

	res, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

type fileUri struct {
	File string
}

func (u fileUri) download() (string, error) {
	return "to be implemented", nil
}

func (u fileUri) upload() (string, error) {
	return "", errors.New("not implemented")
}

func newSecretUri(uri string) (secretUri, error) {

	u, err := url.Parse(uri)
	if err != nil {
		return secretUri{}, err
	}

	ur := secretUri{}
	ur.Keyvault = u.Host
	s := strings.Split(strings.Replace(u.Path, "/secrets/", "", 1), "/")
	ur.Name = s[0]
	ur.Version = ""
	if len(s) == 2 {
		ur.Version = s[1]
	}
	return ur, nil
}

func newFileUri(uri string) (fileUri, error) {
	return fileUri{}, nil
}

func newUri(uri string) (generalUri, error) {

	if strings.HasPrefix(uri, "secret") {
		u, err := newSecretUri(uri)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}

	if strings.HasPrefix(uri, "file") {
		u, err := newFileUri(uri)
		if err != nil {
			return nil, err
		}
		return &u, nil
	}

	return nil, errors.New("invalid full-URL")
}
