package helm_keyvault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	VAULT_NAME_LENGTH              = 24
	VAULT_LOCATION                 = "westeurope"
	KEY_VAULT_ADMINISTRATOR_POLICY = "00482a5a-887f-4fb3-b363-3b7fe8e74483"
)

// skipIntegration - skip integration tests if envrionment isnt specified
func skipIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping integration tests")
	}
}

// randomString - generate a random string
// https://golangdocs.com/generate-random-string-in-golang
func randomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, 8)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// getObjectId - retrieve the object id of the user or service principal running the integration tests
func getObjectId(credentials *azidentity.DefaultAzureCredential) (string, error) {

	t, err := credentials.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes:   []string{"https://graph.microsoft.com/.default"},
		TenantID: "",
	})
	if err != nil {
		return "", err
	}

	// split the received jwt token and decode its payload
	ts := strings.Split(t.Token, ".")
	payload, err := base64.RawStdEncoding.DecodeString(ts[1])
	if err != nil {
		return "", err
	}

	// parse the token and extract the oid field
	var parsed map[string]interface{}
	_ = json.Unmarshal(payload, &parsed)
	if oid, contains := parsed["oid"]; contains {
		return oid.(string), nil
	}
	return "", errors.New("Unable to find OID in token. Aborting")
}

// getArmVaultClient - return an arm client which we can use to create a keyvault resource
func getArmVaultClient(subscription string, credentials *azidentity.DefaultAzureCredential) *armkeyvault.VaultsClient {
	return armkeyvault.NewVaultsClient(subscription, credentials, nil)

}

// getArmSecretsClient - return an arm client which we can use to manage secrets
func getArmSecretsClient(subscription string, credentials *azidentity.DefaultAzureCredential) *armkeyvault.SecretsClient {
	return armkeyvault.NewSecretsClient(subscription, credentials, nil)
}

// getArmKeysClient - return an arm client which we can use to manage secrets
func getArmKeysClient(subscription string, credentials *azidentity.DefaultAzureCredential) *armkeyvault.KeysClient {
	return armkeyvault.NewKeysClient(subscription, credentials, nil)

}

// getKeyVault - check if the keyvault exists. returns nil if it does
func getKeyVault(resourcegroup string, keyvault string, client *armkeyvault.VaultsClient) (armkeyvault.VaultsGetResponse, error) {
	kv, err := client.Get(context.Background(), resourcegroup, keyvault, nil)
	// we assume the keyvault exists if we receive no error (if the keyvault doesnt exist we get a armkeyvault.Clouderror with code ResourceNotFound)
	return kv, err
}

// createKeyVault - create the keyvault
func createKeyVault(resourcegroup string, keyvault string, tenantid string, oid string, vaultsClient *armkeyvault.VaultsClient) error {

	sku := armkeyvault.SKU{
		Family: armkeyvault.SKUFamilyA.ToPtr(),
		Name:   armkeyvault.SKUNameStandard.ToPtr(),
	}

	location := VAULT_LOCATION

	certPermissions := []*armkeyvault.CertificatePermissions{func(p armkeyvault.CertificatePermissions) *armkeyvault.CertificatePermissions { return &p }(armkeyvault.CertificatePermissionsAll)}
	keysPermissions := []*armkeyvault.KeyPermissions{func(p armkeyvault.KeyPermissions) *armkeyvault.KeyPermissions { return &p }(armkeyvault.KeyPermissionsAll)}
	secretsPermissions := []*armkeyvault.SecretPermissions{func(p armkeyvault.SecretPermissions) *armkeyvault.SecretPermissions { return &p }(armkeyvault.SecretPermissionsAll)}

	accessPermissions := armkeyvault.Permissions{
		Certificates: certPermissions,
		Keys:         keysPermissions,
		Secrets:      secretsPermissions,
		Storage:      nil,
	}
	accessPolicy := armkeyvault.AccessPolicyEntry{
		ObjectID:    &oid,
		Permissions: &accessPermissions,
		TenantID:    &tenantid,
		//ApplicationID: nil,
	}
	accessPolicies := []*armkeyvault.AccessPolicyEntry{&accessPolicy}

	properties := armkeyvault.VaultProperties{
		SKU:                     &sku,
		TenantID:                &tenantid,
		AccessPolicies:          accessPolicies,
		CreateMode:              armkeyvault.CreateModeDefault.ToPtr(),
		EnablePurgeProtection:   func(b bool) *bool { return &b }(true),  // convert bool to pointer...
		EnableRbacAuthorization: func(b bool) *bool { return &b }(false), // rbac is recommended but blows up the testing suite complexity
		EnableSoftDelete:        func(b bool) *bool { return &b }(false),
	}

	params := armkeyvault.VaultCreateOrUpdateParameters{
		Location:   &location,
		Properties: &properties,
	}

	poll, err := vaultsClient.BeginCreateOrUpdate(context.Background(), resourcegroup, keyvault, params, nil)
	if err != nil {
		return err
	}

	_, err = poll.PollUntilDone(context.Background(), 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}

// removeKeys() - remove all keys in the keyvault
func removeKeys(keysClient *armkeyvault.KeysClient) error {

	return nil
}

//removeKeyVault() - remove the keyvault
func removeKeyVault(resourcegroup string, keyvault string, vaultsClient *armkeyvault.VaultsClient) error {

	// retrieve and delete all secrets

	// retrieve and delete all keys

	// delete the keyvault
	_, err := vaultsClient.Delete(context.Background(), resourcegroup, keyvault, nil)
	println(err)
	if err != nil {
		return err
	}

	return nil
}

// IntegrationTestSuite - Run keyvault integration tests
// Attention: For the test suite to work the azure identity needs to have permissions
// to create and delete keyvaults in the given resource group.
type IntegrationTestSuite struct {
	suite.Suite

	AzureTenantId      string
	AzureSubscription  string
	AzureResourceGroup string
	AzureKeyVaultName  string
	Credentials        *azidentity.DefaultAzureCredential
	ObjectId           string
	VaultsClient       *armkeyvault.VaultsClient
	KeysClient         *armkeyvault.KeysClient
}

// SetupSuite - Create Keyvault, Make sure
func (s *IntegrationTestSuite) SetupSuite() {
	var err error

	// load azure configuration from env
	log.Info("Parse environment variables")
	s.AzureTenantId = os.Getenv("AZURE_TENANT_ID")
	if s.AzureTenantId == "" {
		log.Fatalf("Please specify the AZURE_TENANT_ID to use for the integration tests.")
	}
	s.AzureSubscription = os.Getenv("AZURE_SUBSCRIPTION")
	if s.AzureSubscription == "" {
		log.Fatalf("Please specify the AZURE_SUBSCRIPTION to use for the integration tests.")
	}
	s.AzureResourceGroup = os.Getenv("AZURE_RESOURCE_GROUP")
	if s.AzureResourceGroup == "" {
		log.Fatalf("Please specify the AZURE_RESOURCE_GROUP to use for the integration tests.")
	}
	s.AzureKeyVaultName = os.Getenv("AZURE_KEYVAULT_NAME")
	if s.AzureKeyVaultName == "" {
		log.Fatalf("Please specify the AZURE_KEYVAULT_NAME to use for the integration tests.")
	}
	// append a random 8 char string to the keyvault name
	// this for safety reasons as the test suite will remove the keyvault at the end of the run
	// just to make sure no kind of race condition can lead to removal of a "real" keyvault
	// make sure the keyvault name wont go over its max size when the random string is appended
	if len(s.AzureKeyVaultName) > (VAULT_NAME_LENGTH - 8) {
		s.AzureKeyVaultName = s.AzureKeyVaultName[0:(VAULT_NAME_LENGTH - 1 - 8)]
	}
	s.AzureKeyVaultName = fmt.Sprintf("%s-%s", s.AzureKeyVaultName, randomString())

	log.Infof("Retrieve azure credentials")
	// setup credentials object to be used with arm
	s.Credentials, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	// retrieve the current users or service principals object id - used to assign policies to the keyvault
	s.ObjectId, err = getObjectId(s.Credentials)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Setup clients")
	s.VaultsClient = getArmVaultClient(s.AzureSubscription, s.Credentials)
	s.KeysClient = getArmKeysClient(s.AzureSubscription, s.Credentials)

	log.Infof("Create keyvault %s in resource group %s (subscription: %s)", s.AzureKeyVaultName, s.AzureResourceGroup, s.AzureSubscription)
	//// check if keyvault exists. if it does abort the operation
	_, err = getKeyVault(s.AzureResourceGroup, s.AzureKeyVaultName, s.VaultsClient)
	if err == nil {
		log.Fatalf("KeyVault %s already exists. Aborting", s.AzureKeyVaultName)
	}

	// create the keyvault
	err = createKeyVault(s.AzureResourceGroup, s.AzureKeyVaultName, s.AzureTenantId, s.ObjectId, s.VaultsClient)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

}

// TearDownSuite - Cleanup azure resources
func (s *IntegrationTestSuite) TearDownSuite() {

	log.Infof("Remove keyvault %s in resource group %s (subscription: %s)", s.AzureKeyVaultName, s.AzureResourceGroup, s.AzureSubscription)

	_, err := getKeyVault(s.AzureResourceGroup, s.AzureKeyVaultName, s.VaultsClient)
	if err != nil {
		// keyvault doesnt exist, nothing to do here
		return
	}

	//kty := armkeyvault.JSONWebKeyTypeRSA
	//props := armkeyvault.KeyProperties{
	//	KeySize: func(s int32) *int32 { return &s }(4096),
	//	Kty:     &kty,
	//}
	//params := armkeyvault.KeyCreateParameters{
	//	Properties: &props,
	//}
	//
	//_, err = s.KeysClient.CreateIfNotExist(context.Background(), s.AzureResourceGroup, s.AzureKeyVaultName, "mykey", params, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//// keyvault exists so lets remove it
	//err = removeKeyVault(s.AzureResourceGroup, s.AzureKeyVaultName, s.VaultsClient)
	//if err != nil {
	//	log.Fatal(err)
	//}

}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *IntegrationTestSuite) TestCreateKey() {

	suite.Equal(suite.AzureResourceGroup, "")
}
func TestExampleTestSuite(t *testing.T) {
	skipIntegration(t)
	suite.Run(t, new(IntegrationTestSuite))
}
