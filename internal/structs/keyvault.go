package structs

import (
	"fmt"
	mskeyvault "github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"net/url"
	"strings"
)

// NewKeyvault - returns a new keyvault struct with valid authorizer and client
// defined as variable to make it easy to override the function inside testting for the cmd package
var NewKeyVault = func(name string) (keyvault.KeyvaultInterface, error) {
	// create keyvault and setup authorizer for it
	kv := keyvault.Keyvault{}
	var err error
	kv.Authorizer, err = kv.NewAuthorizer()
	if err != nil {
		return &keyvault.Keyvault{}, err
	}
	kv.Client = mskeyvault.New()
	kv.Client.Authorizer = kv.Authorizer
	kv.SetKeyvaultName(name)
	return &kv, nil
}

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
