package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
)

// ListKeys - List all secrets in the keyvault
func ListKeys(kv string) error {
	s, err := keyvault.ListKeys(kv)
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

// BackupKey - Backup an azure keyvault key
func BackupKey(kv string, k string, f string) error {
	err := keyvault.BackupKey(kv, k, f)
	return err
}
