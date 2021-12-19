package structs

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/url"
	"strings"
)

//https://<keyvault-name>.vault.azure.net/<type>/<objectname>/<objectversion>"
type KeyvaultObjectId string

// NewKeyVaultObjectId - Return a Keyvault Id
func NewKeyvaultObjectId(kv string, ty string, name string, ver string) KeyvaultObjectId {
	return KeyvaultObjectId(fmt.Sprintf("https://%s.%s/%s/%s/%s", kv, azure.PublicCloud.KeyVaultDNSSuffix, ty, name, ver))
}

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
