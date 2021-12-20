package structs

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"testing"
)

// mocking keyvault until i figure out how to mock the real deal.
type MockKeyvault struct {
	Name string
}

func (m MockKeyvault) EncryptString(key string, version string, encoded string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockKeyvault) DecryptString(key string, version string, encrypted string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockKeyvault) ListKeys() ([]keyvault.KeyBundle, error) {
	var keys []keyvault.KeyBundle

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf(
			"https://%s.%s/keys/%s/%s",
			m.Name,
			azure.PublicCloud.KeyVaultDNSSuffix,
			fmt.Sprintf("key-%v", i),
			"123456789",
		)
		keys = append(keys, keyvault.KeyBundle{
			Key: &keyvault.JSONWebKey{
				Kid: &id,
			},
		})
	}

	return keys, nil
}

func (m MockKeyvault) BackupKey(key string) (string, error) {
	return key, nil
}

func (m MockKeyvault) BackupSecret(secret string) (string, error) {
	return secret, nil
}

func (m MockKeyvault) CreateKey(key string) (keyvault.KeyBundle, error) {
	return keyvault.KeyBundle{}, nil
}

func (m MockKeyvault) GetKey(key string, version string) (keyvault.KeyBundle, error) {

	id := fmt.Sprintf("https://%s.%s/keys/%s/%s", m.Name, azure.PublicCloud.KeyVaultDNSSuffix, key, version)

	return keyvault.KeyBundle{
		Key: &keyvault.JSONWebKey{
			Kid: &id,
		},
	}, nil
}

func (m MockKeyvault) NewAuthorizer() (autorest.Authorizer, error) {
	return nil, nil
}

func (m MockKeyvault) GetKeyvaultName() string {
	return m.Name
}

func (m MockKeyvault) GetSecret(name string, version string) (keyvault.SecretBundle, error) {

	id := fmt.Sprintf("https://%s.%s/secrets/%s/%s", m.Name, azure.PublicCloud.KeyVaultDNSSuffix, name, version)
	value := "My little secret!"

	secret := keyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}

	return secret, nil
}

func (m MockKeyvault) PutSecret(name string, value string) (keyvault.SecretBundle, error) {

	id := fmt.Sprintf("https://%s.%s/secrets/%s/%s", m.Name, azure.PublicCloud.KeyVaultDNSSuffix, name, "myversion")

	secret := keyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}

	return secret, nil
}

func (m MockKeyvault) ListSecrets() ([]keyvault.SecretBundle, error) {

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

func TestNewKeyvaultObjectId(t *testing.T) {
	assert := assert.New(t)

	objectid := NewKeyvaultObjectId("mykeyvault", "mytype", "myname", "myversion")

	assert.Equal("mykeyvault", objectid.GetKeyvault(), "should be equal")
	assert.Equal("mytype", objectid.GetType(), "should be equal")
	assert.Equal("myname", objectid.GetName(), "should be equal")
	assert.Equal("myversion", objectid.GetVersion(), "should be equal")
}

func TestNewKeyvaultObjectIdNoVersion(t *testing.T) {
	assert := assert.New(t)

	objectid := NewKeyvaultObjectId("mykeyvault", "mytype", "myname", "")

	assert.Empty(objectid.GetVersion(), "should be empty")
}
