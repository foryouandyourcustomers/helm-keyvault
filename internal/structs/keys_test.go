package structs

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/stretchr/testify/assert"
	"testing"
)

// mocking keyvault until i figure out how to mock the real deal.
type MockKeysKeyvault struct {
	Name string
}

func (m MockKeysKeyvault) EncryptString(key string, version string, encoded string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockKeysKeyvault) DecryptString(key string, version string, encrypted string) (keyvault.KeyOperationResult, error) {
	return keyvault.KeyOperationResult{}, nil
}

func (m MockKeysKeyvault) ListKeys() ([]keyvault.KeyBundle, error) {
	return nil, nil
}

func (m MockKeysKeyvault) BackupKey(key string) (string, error) {
	return "", nil
}

func (m MockKeysKeyvault) CreateKey(key string) (keyvault.KeyBundle, error) {
	return keyvault.KeyBundle{}, nil
}

func (m MockKeysKeyvault) GetKey(key string, version string) (keyvault.KeyBundle, error) {
	return keyvault.KeyBundle{}, nil
}

func (m MockKeysKeyvault) NewAuthorizer() (autorest.Authorizer, error) {
	return nil, nil
}

func (m MockKeysKeyvault) GetKeyvaultName() string {
	return m.Name
}

func (m MockKeysKeyvault) GetSecret(name string, version string) (keyvault.SecretBundle, error) {
	return keyvault.SecretBundle{}, nil
}

func (m MockKeysKeyvault) PutSecret(name string, value string) (keyvault.SecretBundle, error) {

	return keyvault.SecretBundle{}, nil
}

func (m MockKeysKeyvault) ListSecrets() ([]keyvault.SecretBundle, error) {
	return nil, nil
}

func TestNewKey(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(nil, "tbd")
}

func TestKey_Backup(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(nil, "tbd")
}

func TestKey_Get(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(nil, "tbd")
}

func TestKey_Create(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(nil, "tbd")
}

func TestKeyList_List(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(nil, "tbd")
}
