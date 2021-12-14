package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"io/ioutil"
	"net/url"
	"strings"
)

type secretUri struct {
	Keyvault string
	Name     string
	Version  string
}

func newSecretUri(uri string) (*secretUri, error) {
	// parse the given uri
	// keyvault://<keyvaultname>/secrets/<secretname>/<version>
	u, err := url.Parse(uri)
	if err != nil {
		return &secretUri{}, err
	}

	// retrieve keyvault, secret and optional secret version from parsed uri
	ur := secretUri{}
	ur.Keyvault = u.Host
	s := strings.Split(strings.Replace(u.Path, "/secrets/", "", 1), "/")
	ur.Name = s[0]
	ur.Version = ""
	if len(s) == 2 {
		ur.Version = s[1]
	}
	return &ur, nil
}

// GetSecret - Get secret from given keyvault
func GetSecret(id string) error {
	err := DownloadSecret(id)
	return err
}

// PutSecret - Encode file and put secret into keyvault
func PutSecret(id string, f string) error {

	// parse given id into its keyvault components
	su, err := newSecretUri(id)

	if err != nil {
		return err
	}

	// read file and convert it to base64
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	e := base64.StdEncoding.EncodeToString(c)

	s, err := keyvault.PutSecret(su.Keyvault, su.Name, e)
	if err != nil {
		return err
	}

	res, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

// DownloadSecret - Download and decode secret to be used as downloader plugin
func DownloadSecret(uri string) error {
	su, err := newSecretUri(uri)
	if err != nil {
		return err
	}

	// get secret from keyvault and print it
	secret, err := keyvault.GetSecret(su.Keyvault, su.Name, su.Version)
	if err != nil {
		return err
	}

	// decode secret and print
	dec, err := secret.Decode()
	if err != nil {
		return err
	}
	fmt.Print(dec)
	return nil
}
