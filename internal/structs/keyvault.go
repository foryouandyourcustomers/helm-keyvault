package structs

import (
	"net/url"
	"strings"
)

//https://<keyvault-name>.vault.azure.net/<type>/<objectname>/<objectversion>"
type KeyvaultObjectId string

// GetKeyvault - Get the keyvault name from the ObjectId
func (k *KeyvaultObjectId) GetKeyvault() string {
	kv, _ := url.Parse(string(*k))
	h := strings.Split(kv.Host, ".")
	return h[0]
}

func (k *KeyvaultObjectId) GetType() string {
	kv, _ := url.Parse(string(*k))
	p, _ := splitPath(kv.Path)
	return p[1]
}

func (k *KeyvaultObjectId) GetName() string {
	kv, _ := url.Parse(string(*k))
	p, _ := splitPath(kv.Path)
	return p[2]
}

func (k *KeyvaultObjectId) GetVersion() string {
	kv, _ := url.Parse(string(*k))
	p, _ := splitPath(kv.Path)
	if len(p) < 4 {
		return ""
	}
	return p[3]
}
