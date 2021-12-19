package cmd

import (
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"strings"
	"time"
)

// EncryptFile - encrypt the given file with the given key
func EncryptFile(kv string, k string, v string, f string) error {

	// setup encrypted file object
	ef := structs.EncryptedFile{}

	// add values to encoded file struct
	ef.Kid = structs.CreateKeyVaultId(kv, "keys", k, v)

	// load file
	var err error
	ef.EncodedData, err = ef.LoadFile(f)
	if err != nil {
		return err
	}

	// encrypt the given dats
	ef.EncryptedData, err = ef.EncryptData()
	if err != nil {
		return err
	}
	ef.LastModified = structs.JTime(time.Now())

	// write file
	err = ef.WriteEncryptedFile(f)
	return err
}

// DecryptFile - decrypt the given file with the key specified in the encrypted
// file. The keyvault and namespace can be overwritten via paraeters/env vars
func DecryptFile(kv string, k string, v string, f string) error {

	ef := structs.EncryptedFile{}
	ef, err := ef.LoadEncryptedFile(f)
	if err != nil {
		return err
	}

	// decrypt data, overwrite given kid with optional key
	ef.EncodedData, err = ef.DecryptData(kv, k, v)
	if err != nil {
		return err
	}

	// write decrypted data to disk
	fn := strings.Replace(f, ".enc", "", 1)
	err = ef.WriteFile(fn)
	return err

}
