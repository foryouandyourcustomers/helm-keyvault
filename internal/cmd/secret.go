package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"io/ioutil"
)

// GetSecret - Get secret from given keyvault
func GetSecret(kv string, sn string, ve string) error {

	// retrieve and decode base64 encoded secret
	su, err := keyvault.GetSecret(kv, sn, ve)
	if err != nil {
		return err
	}

	// retrieve and decode secret content
	value, err := su.Decode()
	if err != nil {
		return err
	}
	fmt.Print(value)
	return nil
}

// PutSecret - Encode file and put secret into keyvault
func PutSecret(kv string, sn string, f string) error {

	// read file and convert it to base64
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	e := base64.StdEncoding.EncodeToString(c)

	s, err := keyvault.PutSecret(kv, sn, e)
	if err != nil {
		return err
	}

	res, err := json.Marshal(s)
	if err != nil {
		return err
	}

	fmt.Print(string(res))
	return nil
}

// ListSecrets - List all secrets in the keyvault
func ListSecrets(kv string) error {
	s, err := keyvault.ListSecrets(kv)
	if err != nil {
		return err
	}
	j, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fmt.Print(string(j))
	return nil
}
