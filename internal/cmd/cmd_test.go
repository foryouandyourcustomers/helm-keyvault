package cmd

import (
	"fmt"
	mskeyvault "github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
)

func newMockKeyVault(name string) (keyvault.KeyvaultInterface, error) {
	kv := MockKeyVault{}
	kv.SetKeyvaultName(name)

	return &kv, nil
}

type MockKeyVault struct {
	Name    string
	BaseUrl string
}

func (m *MockKeyVault) NewAuthorizer() (autorest.Authorizer, error) {
	return nil, nil
}

func (m *MockKeyVault) GetKeyvaultName() string {
	return m.Name
}

func (m *MockKeyVault) SetKeyvaultName(name string) {
	if name != "" {
		m.Name = name
		m.BaseUrl = fmt.Sprintf("https://%s.%s", name, azure.PublicCloud.KeyVaultDNSSuffix)
	}
}

func (m *MockKeyVault) GetSecret(sn string, sv string) (mskeyvault.SecretBundle, error) {
	id := string(structs.NewKeyvaultObjectId(m.Name, "secrets", sn, sv))
	value := "Exammple Value"
	return mskeyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}, nil
}

func (m *MockKeyVault) PutSecret(name string, value string) (mskeyvault.SecretBundle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) ListSecrets() ([]mskeyvault.SecretBundle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) BackupSecret(sn string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) EncryptString(key string, version string, encoded string) (mskeyvault.KeyOperationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) DecryptString(key string, version string, encrypted string) (mskeyvault.KeyOperationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) ListKeys() ([]mskeyvault.KeyBundle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) BackupKey(key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) CreateKey(key string) (mskeyvault.KeyBundle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockKeyVault) GetKey(key string, version string) (mskeyvault.KeyBundle, error) {
	//TODO implement me
	panic("implement me")
}
