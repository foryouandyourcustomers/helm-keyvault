package structs

import (
	"errors"
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
	"strings"
)

func splitPath(p string) ([]string, error) {
	paths := strings.Split(p, "/")
	if len(paths) != 4 {
		return []string{}, errors.New(fmt.Sprintf("Invalid keyvault path '%s'", p))
	}
	return paths, nil
}

func CreateKeyVaultId(kv string, ty string, name string, ver string) KeyvaultObjectId {
	return KeyvaultObjectId(fmt.Sprintf("https://%s.%s/%s/%s/%s", kv, azure.PublicCloud.KeyVaultDNSSuffix, ty, name, ver))
}
