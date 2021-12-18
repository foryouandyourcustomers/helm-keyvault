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

func (k *Key) Create() error {
	// first check if the key already exists
	kb, err := keyvault.GetKey(k.KeyVault, k.Name, k.Version)
	// abort here if the key can be retrieved

	if err == nil {
		return errors.New("Key already exists.")
	}

	kb, err = keyvault.CreateKey(k.KeyVault, k.Name)
	if err != nil {
		return err
	}

	koid := KeyvaultObjectId(*kb.Key.Kid)
	kn, _ := koid.GetName()
	kkv, _ := koid.GetKeyvault()
	kve, _ := koid.GetKeyvault()

	k.Kid = koid
	k.Name = kn
	k.Version = kve
	k.KeyVault = kkv

	return nil
}

type KeyList struct {
	Keys []Key `json:"keys,omitempty"`
}

func (sl *KeyList) List(kv string) error {
	sk, err := keyvault.ListKeys(kv)
	if err != nil {
		return err
	}

	for _, k := range sk {
		koid := KeyvaultObjectId(*k.Key.Kid)
		sn, _ := koid.GetName()
		skv, _ := koid.GetKeyvault()
		sve, _ := koid.GetVersion()

		sl.Keys = append(sl.Keys,
			Key{
				Kid:      koid,
				Name:     sn,
				KeyVault: skv,
				Version:  sve,
			})
	}

	return nil
}
