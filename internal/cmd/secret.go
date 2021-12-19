package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"io/ioutil"
)

// GetSecret - Get secret from given keyvault
func GetSecret(kv string, sn string, ve string) error {

	// retrieve and decode base64 encoded secret
	sec := structs.Secret{
		Name:     sn,
		Version:  ve,
		KeyVault: kv,
	}

	sec, err := sec.Get()
	if err != nil {
		return err
	}

	j, err := json.Marshal(sec)
	if err != nil {
		return err
	}
	fmt.Print(string(j))
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

	// put secret to keyvault
	sec := structs.Secret{
		Name:     sn,
		KeyVault: kv,
		Value:    e,
	}

	sec, err = sec.Put()
	if err != nil {
		return err
	}
	j, err := json.Marshal(sec)
	fmt.Print(string(j))
	return nil
}

// ListSecrets - List all secrets in the keyvault
func ListSecrets(kv string) error {

	// initialize list
	sl := structs.SecretList{}

	var err error
	sl.Secrets, err = sl.List(kv)
	if err != nil {
		return err
	}

	j, err := json.Marshal(sl)
	fmt.Print(string(j))
	return nil
}
