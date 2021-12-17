package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

// EncryptFile - encrypt the given file with the given key
func EncryptFile(kv string, k string, v string, f string) error {

	// load file and encode is as base64 url (non padded)
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	e := base64.URLEncoding.EncodeToString(c)

	// encrypt file content with keyvault
	s, err := keyvault.EncryptString(kv, k, v, e)
	if err != nil {
		return err
	}

	// write file
	ef := structs.EncryptedFile{
		Kid:          structs.KeyvaultObjectId{*s.Kid},
		Data:         *s.Result,
		LastModified: structs.JTime(time.Now()),
	}

	efj, err := json.MarshalIndent(ef, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s.enc", f), efj, 0644)
	if err != nil {
		return err
	}
	return nil
}

// DecryptFile - decrypt the given file with the key specified in the encrypted
// file. The keyvault and namespace can be overwritten via paraeters/env vars
func DecryptFile(kv string, k string, v string, f string) error {

	// load file and encode is as base64 url (non padded)
	c, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	con := structs.EncryptedFile{}
	err = json.Unmarshal(c, &con)
	if err != nil {
		return err
	}

	// check if overwrite values are set for keyvault,
	// key and key version
	kvname := kv
	if kvname == "" {
		kvname, err = con.Kid.GetKeyvault()
		if err != nil {
			return err
		}
	}
	key := k
	if key == "" {
		key, err = con.Kid.GetName()
		if err != nil {
			return err
		}
	}
	version := v
	if version == "" {
		version, err = con.Kid.GetVersion()
		if err != nil {
			return err
		}
	}

	// decrypt file contents and write file
	dec, err := keyvault.DecryptString(kvname, key, version, con.Data)
	if err != nil {
		return err
	}
	// decode base64 value
	val, err := base64.URLEncoding.DecodeString(*dec.Result)
	if err != nil {
		return err
	}

	// write value to file
	fn := strings.Replace(f, ".enc", "", 1)
	err = os.WriteFile(fn, val, 0644)
	if err != nil {
		return err
	}
	return nil
}
