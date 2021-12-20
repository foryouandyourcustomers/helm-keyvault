package structs

import (
	"encoding/base64"
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSecret(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}

	secret := NewSecret(mock, "mysecret", "myversion")

	assert.Empty(secret.Value, "should be empty")
	assert.Equal(KeyvaultObjectId(fmt.Sprintf("https://mykeyvault.%s/secrets/mysecret/myversion", azure.PublicCloud.KeyVaultDNSSuffix)), secret.Id, "should be equal")
	assert.Equal("mykeyvault", secret.KeyVault.GetKeyvaultName(), "should be equal")
	assert.Equal("mysecret", secret.Name, "should be equal")
	assert.Equal("myversion", secret.Version, "should be equal")
}

func TestSecret_Get(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	secret := NewSecret(mock, "mysecret", "myversion")

	s, err := secret.Get()
	assert.Nil(err, "should be nil")
	assert.Equal(string(s.Id), fmt.Sprintf("https://%s.%s/secrets/%s/%s", "mykeyvault", azure.PublicCloud.KeyVaultDNSSuffix, "mysecret", "myversion"))
	assert.Equal(s.Name, "mysecret", "should be equal")
	assert.Equal(s.Value, "My little secret!", "should be equal")
	assert.Equal(s.Version, "myversion", "should be equal")
}

func TestSecret_Put(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	secret := NewSecret(mock, "mysecret", "")
	secret.Value = "My little secret!"
	s, err := secret.Put()

	assert.Nil(err, "should be nil")
	assert.Equal(string(s.Id), fmt.Sprintf("https://%s.%s/secrets/%s/%s", "mykeyvault", azure.PublicCloud.KeyVaultDNSSuffix, "mysecret", "myversion"))
	assert.Equal(s.Name, "mysecret", "should be equal")
	assert.Equal(s.Value, "My little secret!", "should be equal")
	assert.Equal(s.Version, "myversion", "should be equal")
}

func TestSecret_Decode(t *testing.T) {
	assert := assert.New(t)

	rawstring := "My little secret!"
	encoded_string := base64.StdEncoding.EncodeToString([]byte(rawstring))

	mock := MockKeyvault{Name: "mykeyvault"}
	secret := NewSecret(mock, "mysecret", "")
	secret.Value = encoded_string
	secret, _ = secret.Put()

	dec, err := secret.Decode()
	assert.Nil(err, "should be nil")
	assert.Equal(dec, "My little secret!")

	secret.Value = rawstring
	secret, _ = secret.Put()

	dec, err = secret.Decode()
	assert.Empty(dec, "should be empy")
	assert.Error(err, "should be error")
}

func TestSecretList_List(t *testing.T) {

	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	sl := SecretList{}
	sl.Secrets, _ = sl.List(&mock)

	assert.Len(sl.Secrets, 5, "should be 5")
}
