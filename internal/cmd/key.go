package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
)

// ListKeys - List all secrets in the keyvault
func ListKeys(kv string) error {

	// initialize keyvault object
	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	// initialize list
	sl := structs.KeyList{}

	sl.Keys, err = sl.List(&keyvault)
	if err != nil {
		return err
	}

	j, err := json.Marshal(sl)
	fmt.Print(string(j))
	return nil
}

// BackupKey - Backup an azure keyvault key
func BackupKey(kv string, k string, f string) error {

	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	key := structs.NewKey(&keyvault, k, "")
	err = key.Backup(f)
	return err
}

// CreateKey - Create an azure keyvault key
func CreateKey(kv string, k string) error {

	keyvault, err := structs.NewKeyvault(kv)
	if err != nil {
		return err
	}

	key := structs.NewKey(&keyvault, k, "")

	key, err = key.Create()
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
