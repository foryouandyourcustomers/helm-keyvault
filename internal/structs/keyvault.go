package structs

import (
	"net/url"
	"strings"
)

//https://<keyvault-name>.vault.azure.net/<type>/<objectname>/<objectversion>"
type KeyvaultObjectId string

// GetKeyvault - Get the keyvault name from the ObjectId
func (k *KeyvaultObjectId) GetKeyvault() (string, error) {
	kv, err := url.Parse(string(*k))
	if err != nil {
		return "", err
	}

	h := strings.Split(kv.Host, ".")
	return h[0], nil
}

func (k *KeyvaultObjectId) GetType() (string, error) {
	kv, err := url.Parse(string(*k))
	if err != nil {
		return "", err
	}

	p, err := splitPath(kv.Path)
	if err != nil {
		return "", err
	}

	return p[1], nil
}

func (k *KeyvaultObjectId) GetName() (string, error) {
	kv, err := url.Parse(string(*k))
	if err != nil {
		return "", err
	}

	p, err := splitPath(kv.Path)
	if err != nil {
		return "", err
	}

	return p[2], nil
}

func (k *KeyvaultObjectId) GetVersion() (string, error) {
	kv, err := url.Parse(string(*k))
	if err != nil {
		return "", err
	}

	p, err := splitPath(kv.Path)
	if err != nil {
		return "", err
	}

	if len(p) < 4 {
		return "", nil
	}

	return p[3], nil
}
