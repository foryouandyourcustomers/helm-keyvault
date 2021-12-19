package structs

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSecret(t *testing.T) {
	assert := assert.New(t)

	secret := NewSecret("mykeyvault", "mysecret", "myversion")

	assert.Empty(secret.Value, "should be empty")
	assert.Equal(KeyvaultObjectId(fmt.Sprintf("https://mykeyvault.%s/secrets/mysecret/myversion", azure.PublicCloud.KeyVaultDNSSuffix)), secret.Id, "should be equal")
	assert.Equal("mykeyvault", secret.KeyVault, "should be equal")
	assert.Equal("mysecret", secret.Name, "should be equal")
	assert.Equal("myversion", secret.Version, "should be equal")
}

func TestSecret_Get(t *testing.T) {
	assert := assert.New(t)

	secret := NewSecret("mykeyvault", "mysecret", "myversion")

	s, err := secret.Get()

	assert.Nil(err, "should be nil")

	println(s.Value)

}
