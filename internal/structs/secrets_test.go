package structs

import (
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"testing"
)

// mocking keyvault until i figure out how to mock the real deal.
type MockSecretsKeyvault struct {
	Name string
}

func (m MockSecretsKeyvault) EncryptString(key string, version string, encoded string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockSecretsKeyvault) DecryptString(key string, version string, encrypted string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockSecretsKeyvault) ListKeys() ([]keyvault.KeyBundle, error) {
	return nil, nil
}

func (m MockSecretsKeyvault) BackupKey(key string) (string, error) {
	return "", nil
}

func (m MockSecretsKeyvault) CreateKey(key string) (keyvault.KeyBundle, error) {
	return keyvault.KeyBundle{}, nil
}

func (m MockSecretsKeyvault) GetKey(key string, version string) (keyvault.KeyBundle, error) {
	return keyvault.KeyBundle{}, nil
}

func (m MockSecretsKeyvault) NewAuthorizer() (autorest.Authorizer, error) {
	return nil, nil
}

func (m MockSecretsKeyvault) GetKeyvaultName() string {
	return m.Name
}

func (m MockSecretsKeyvault) GetSecret(name string, version string) (keyvault.SecretBundle, error) {

	id := fmt.Sprintf("https://%s.%s/secrets/%s/%s", m.Name, azure.PublicCloud.KeyVaultDNSSuffix, name, version)
	value := "My little secret!"

	secret := keyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}

	return secret, nil
}

func (m MockSecretsKeyvault) PutSecret(name string, value string) (keyvault.SecretBundle, error) {

	id := fmt.Sprintf("https://%s.%s/secrets/%s/%s", m.Name, azure.PublicCloud.KeyVaultDNSSuffix, name, "myversion")

	secret := keyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}

	return secret, nil
}

func (m MockSecretsKeyvault) ListSecrets() ([]keyvault.SecretBundle, error) {

	var secrets []keyvault.SecretBundle

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf(
			"https://%s.%s/secrets/%s/%s",
			m.Name,
			azure.PublicCloud.KeyVaultDNSSuffix,
			fmt.Sprintf("secret-%v", i),
			"123456789",
		)
		val := fmt.Sprintf("My N-th (%v) secret", i)
		secrets = append(secrets, keyvault.SecretBundle{
			ID:    &id,
			Value: &val,
		})
	}

	return secrets, nil
}

func TestNewSecret(t *testing.T) {
	assert := assert.New(t)

	mock := MockSecretsKeyvault{Name: "mykeyvault"}

	secret := NewSecret(mock, "mysecret", "myversion")

	assert.Empty(secret.Value, "should be empty")
	assert.Equal(KeyvaultObjectId(fmt.Sprintf("https://mykeyvault.%s/secrets/mysecret/myversion", azure.PublicCloud.KeyVaultDNSSuffix)), secret.Id, "should be equal")
	assert.Equal("mykeyvault", secret.KeyVault.GetKeyvaultName(), "should be equal")
	assert.Equal("mysecret", secret.Name, "should be equal")
	assert.Equal("myversion", secret.Version, "should be equal")
}

func TestSecret_Get(t *testing.T) {
	assert := assert.New(t)

	mock := MockSecretsKeyvault{Name: "mykeyvault"}
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

	mock := MockSecretsKeyvault{Name: "mykeyvault"}
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

	mock := MockSecretsKeyvault{Name: "mykeyvault"}
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

	mock := MockSecretsKeyvault{Name: "mykeyvault"}
	sl := SecretList{}
	sl.Secrets, _ = sl.List(&mock)

	assert.Len(sl.Secrets, 5, "should be 5")
}
