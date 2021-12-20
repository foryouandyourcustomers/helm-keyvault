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
	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	sec := structs.NewSecret(&keyvault, sn, ve)

	sec, err = sec.Get()
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
	// retrieve and decode base64 encoded secret
	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	sec := structs.NewSecret(&keyvault, sn, "")
	sec.Value = e

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

	// initialize keyvault object
	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	// inialize secret list
	sl := structs.SecretList{}

	sl.Secrets, err = sl.List(&keyvault)
	if err != nil {
		return err
	}

	j, err := json.Marshal(sl)
	fmt.Print(string(j))
	return nil
}

// BackupSecret - Create a backup of the specified secret
func BackupSecret(kv string, secret string, file string) error {

	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	sec := structs.NewSecret(&keyvault, secret, "")
	err = sec.Backup(file)
	return err
}
