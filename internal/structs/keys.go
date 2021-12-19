package structs

import (
	"errors"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"os"
)

type Key struct {
	Kid      KeyvaultObjectId `json:"kid,omitempty"`
	Name     string           `json:"name,omitempty"`
	KeyVault string           `json:"keyvault,omitempty"`
	Version  string           `json:"version,omitempty"`
}

// Backup - create backup of key and write it into the given file
func (k *Key) Backup(f string) error {
	backup, err := keyvault.BackupKey(k.KeyVault, k.Name)
	if err != nil {
		return err
	}

	fp, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.WriteString(backup)
	if err != nil {
		return err
	}
	return nil
}

// Get - Retrieve key information from keyvault
func (k *Key) Get() (Key, error) {
	kb, err := keyvault.GetKey(k.KeyVault, k.Name, k.Version)
	if err != nil {
		return Key{}, err
	}

	koid := KeyvaultObjectId(*kb.Key.Kid)
	return Key{
		Kid:      koid,
		Name:     koid.GetName(),
		KeyVault: koid.GetKeyvault(),
		Version:  koid.GetVersion(),
	}, nil
}

func (k *Key) Create() (Key, error) {
	// first check if the key already exists
	kb, err := keyvault.GetKey(k.KeyVault, k.Name, k.Version)
	// abort here if the key can be retrieved
	if err == nil {
		return Key{}, errors.New("Key already exists.")
	}

	kb, err = keyvault.CreateKey(k.KeyVault, k.Name)
	if err != nil {
		return Key{}, err
	}

	koid := KeyvaultObjectId(*kb.Key.Kid)
	return Key{
		Kid:      koid,
		Name:     koid.GetName(),
		KeyVault: koid.GetKeyvault(),
		Version:  koid.GetVersion(),
	}, nil

}

type KeyList struct {
	Keys []Key `json:"keys,omitempty"`
}

func (sl *KeyList) List(kv string) ([]Key, error) {
	sk, err := keyvault.ListKeys(kv)
	if err != nil {
		return nil, err
	}

	var keys []Key
	for _, k := range sk {
		koid := KeyvaultObjectId(*k.Key.Kid)
		sn := koid.GetName()
		skv := koid.GetKeyvault()
		sve := koid.GetVersion()

		keys = append(keys,
			Key{
				Kid:      koid,
				Name:     sn,
				KeyVault: skv,
				Version:  sve,
			})
	}

	return keys, nil
}
