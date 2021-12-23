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

func (m *MockKeyVault) GetSecret(name string, version string) (mskeyvault.SecretBundle, error) {
	id := string(structs.NewKeyvaultObjectId(m.Name, "secrets", name, version))
	value := "Exammple Value"
	return mskeyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}, nil
}

func (m *MockKeyVault) PutSecret(name string, value string) (mskeyvault.SecretBundle, error) {
	version := "123456"
	id := string(structs.NewKeyvaultObjectId(m.Name, "secrets", name, version))
	return mskeyvault.SecretBundle{
		ID:    &id,
		Value: &value,
	}, nil
}

func (m *MockKeyVault) ListSecrets() ([]mskeyvault.SecretBundle, error) {

	var secrets []mskeyvault.SecretBundle

	for i := 0; i < 5; i++ {
		id := fmt.Sprintf(
			"https://%s.%s/secrets/%s/%s",
			m.Name,
			azure.PublicCloud.KeyVaultDNSSuffix,
			fmt.Sprintf("secret-%v", i),
			"123456789",
		)
		val := fmt.Sprintf("My N-th (%v) secret", i)
		secrets = append(secrets, mskeyvault.SecretBundle{
			ID:    &id,
			Value: &val,
		})
	}

	return secrets, nil
}
func (m *MockKeyVault) BackupSecret(secret string) (string, error) {
	return secret, nil
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
	return mskeyvault.KeyBundle{}, nil
}

func (m *MockKeyVault) GetKey(key string, version string) (mskeyvault.KeyBundle, error) {
	//TODO implement me
	panic("implement me")
}
