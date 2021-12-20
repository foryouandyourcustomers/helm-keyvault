package structs

import (
	"errors"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"os"
)

// NewKey - create a new Key struct
func NewKey(kv keyvault.KeyvaultInterface, key string, version string) Key {
	return Key{
		Kid:      NewKeyvaultObjectId(kv.GetKeyvaultName(), "keys", key, version),
		Name:     key,
		KeyVault: kv,
		Version:  version,
	}
}

type Key struct {
	Kid      KeyvaultObjectId           `json:"kid,omitempty"`
	Name     string                     `json:"name,omitempty"`
	KeyVault keyvault.KeyvaultInterface `json:"keyvault,omitempty"`
	Version  string                     `json:"version,omitempty"`
}

// Backup - create backup of key and write it into the given file
func (k *Key) Backup(f string) error {
	backup, err := k.KeyVault.BackupKey(k.Name)
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
	kb, err := k.KeyVault.GetKey(k.Name, k.Version)
	if err != nil {
		return Key{}, err
	}

	koid := KeyvaultObjectId(*kb.Key.Kid)
	return Key{
		Kid:      koid,
		Name:     koid.GetName(),
		KeyVault: k.KeyVault,
		Version:  koid.GetVersion(),
	}, nil
}

func (k *Key) Create() (Key, error) {
	// first check if the key already exists
	kb, err := k.KeyVault.GetKey(k.Name, k.Version)
	// abort here if the key can be retrieved
	if err == nil {
		return Key{}, errors.New("Key already exists.")
	}

	kb, err = k.KeyVault.CreateKey(k.Name)
	if err != nil {
		return Key{}, err
	}

	koid := KeyvaultObjectId(*kb.Key.Kid)
	return Key{
		Kid:      koid,
		Name:     koid.GetName(),
		KeyVault: k.KeyVault,
		Version:  koid.GetVersion(),
	}, nil

}

type KeyList struct {
	Keys []Key `json:"keys,omitempty"`
}

func (sl *KeyList) List(kv keyvault.KeyvaultInterface) ([]Key, error) {
	sk, err := kv.ListKeys()
	if err != nil {
		return nil, err
	}

	var keys []Key
	for _, k := range sk {
		koid := KeyvaultObjectId(*k.Key.Kid)
		sn := koid.GetName()
		sve := koid.GetVersion()

		keys = append(keys,
			Key{
				Kid:      koid,
				Name:     sn,
				KeyVault: kv,
				Version:  sve,
			})
	}

	return keys, nil
}
