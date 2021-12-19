package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
)

// ListKeys - List all secrets in the keyvault
func ListKeys(kv string) error {

	// initialize list
	sl := structs.KeyList{}

	var err error
	sl.Keys, err = sl.List(kv)
	if err != nil {
		return err
	}

	j, err := json.Marshal(sl)
	fmt.Print(string(j))
	return nil
}

// BackupKey - Backup an azure keyvault key
func BackupKey(kv string, k string, f string) error {

	key := structs.Key{
		KeyVault: kv,
		Name:     k,
	}

	err := key.Backup(f)
	return err
}

// CreateKey - Create an azure keyvault key
func CreateKey(kv string, k string) error {

	key := structs.Key{
		KeyVault: kv,
		Name:     k,
	}

	key, err := key.Create()
	if err != nil {
		return err
	}

	j, err := json.Marshal(key)
	if err != nil {
		return err
	}
	fmt.Print(string(j))
	return nil

}
