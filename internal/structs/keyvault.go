package structs

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func splitPath(p string) ([]string, error) {
	paths := strings.Split(p, "/")
	if len(paths) != 4 {
		return []string{}, errors.New(fmt.Sprintf("Invalid keyvault path '%s'", p))
	}
	return paths, nil
}

type KeyvaultObjectId struct {
	//https://<keyvault-name>.vault.azure.net/<type>/<objectname>/<objectversion>"
	Id string
}

// GetKeyvault - Get the keyvault name from the ObjectId
func (k *KeyvaultObjectId) GetKeyvault() (string, error) {
	kv, err := url.Parse(k.Id)
	if err != nil {
		return "", err
	}

	h := strings.Split(kv.Host, ".")
	return h[0], nil
}

func (k *KeyvaultObjectId) GetType() (string, error) {
	kv, err := url.Parse(k.Id)
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
	kv, err := url.Parse(k.Id)
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
	kv, err := url.Parse(k.Id)
	if err != nil {
		return "", err
	}

	p, err := splitPath(kv.Path)
	if err != nil {
		return "", err
	}

	return p[3], nil
}
