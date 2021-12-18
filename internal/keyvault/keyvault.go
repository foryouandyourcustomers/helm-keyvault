package keyvault

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/sirupsen/logrus"
	"path"
)

var (
	authorizer autorest.Authorizer
)

const (
	KeyType keyvault.JSONWebKeyType                = "RSA"
	KeySize int32                                  = 4096
	KeyAlgo keyvault.JSONWebKeyEncryptionAlgorithm = keyvault.RSA15
)

// initialize keyvault authorizer
func init() {
	// first try to get authorizer from cli
	var err error

	log.Debug("Try to get authentication from file")
	authorizer, err = kvauth.NewAuthorizerFromFile()
	if err != nil {
		log.Debug("Try to get credentials from envrionment")
		authorizer, err = kvauth.NewAuthorizerFromEnvironment()
		if err != nil {
			log.Debug("Get login info from azure cli")
			authorizer, err = kvauth.NewAuthorizerFromCLI()
			if err != nil {
				panic("Unable to authenticate with AUTH file, ENV vars and local cli. Aborting.")
			}
		}
	}
}

// GetSecret - return a secret object
func GetSecret(kv string, sn string, sv string) (keyvault.SecretBundle, error) {

	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	s, err := c.GetSecret(context.Background(), baseurl, sn, sv)
	if err != nil {
		return keyvault.SecretBundle{}, err
	}

	return s, nil
}

// PutSecret - put secret into keyvault
func PutSecret(kv string, sn string, cn string) (keyvault.SecretBundle, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	ct := "base64"
	sp := keyvault.SecretSetParameters{
		Value:       &cn,
		ContentType: &ct,
	}

	s, err := c.SetSecret(context.Background(), baseurl, sn, sp)
	if err != nil {
		return keyvault.SecretBundle{}, err
	}
	return s, nil
}

// ListSecrets - list all secrets in the specified keyvault
func ListSecrets(kv string) ([]keyvault.SecretBundle, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	ctx := context.Background()
	siter, err := c.GetSecretsComplete(ctx, baseurl, nil)
	if err != nil {
		log.Fatalf("unable to get list of secrets: %v\n", err)
	}

	var s []keyvault.SecretBundle

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.ID)
		b, err := c.GetSecret(context.Background(), baseurl, key, "")
		if err != nil {
			return []keyvault.SecretBundle{}, err
		}

		s = append(s, b)
		err = siter.NextWithContext(ctx)
		if err != nil {
			return []keyvault.SecretBundle{}, err
		}
	}

	return s, nil
}

// EncryptString - encrypt a given file
func EncryptString(kv string, k string, v string, e string) (keyvault.KeyOperationResult, error) {

	// prepare keyvault
	c := keyvault.New()
	c.Authorizer = authorizer
	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)
	ctx := context.Background()
	param := keyvault.KeyOperationsParameters{
		Algorithm: KeyAlgo,
		Value:     &e,
	}
	r, err := c.Encrypt(ctx, baseurl, k, v, param)
	if err != nil {
		return keyvault.KeyOperationResult{}, err
	}

	return r, nil
}

func DecryptString(kv string, k string, v string, e string) (keyvault.KeyOperationResult, error) {

	// prepare keyvault
	c := keyvault.New()
	c.Authorizer = authorizer
	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)
	ctx := context.Background()
	param := keyvault.KeyOperationsParameters{
		Algorithm: KeyAlgo,
		Value:     &e,
	}
	r, err := c.Decrypt(ctx, baseurl, k, v, param)
	if err != nil {
		return keyvault.KeyOperationResult{}, err
	}

	return r, nil
}

// ListKeys - list all keys in the specified keyvault
func ListKeys(kv string) ([]keyvault.KeyBundle, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	ctx := context.Background()
	siter, err := c.GetKeysComplete(ctx, baseurl, nil)
	if err != nil {
		log.Fatalf("unable to get list of keys: %v\n", err)
	}

	var k []keyvault.KeyBundle

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.Kid)
		b, err := c.GetKey(context.Background(), baseurl, key, "")
		if err != nil {
			return []keyvault.KeyBundle{}, err
		}

		k = append(k, b)
		err = siter.NextWithContext(ctx)
		if err != nil {
			return []keyvault.KeyBundle{}, err
		}
	}

	return k, nil
}

// BackupKey - Create a backup of a key which can be used for restoring
func BackupKey(kv string, key string) (string, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)
	kb, err := c.BackupKey(context.Background(), baseurl, key)
	if err != nil {
		return "", err
	}

	dec, err := base64.RawURLEncoding.DecodeString(*kb.Value)
	if err != nil {
		return "", err
	}
	return string(dec), nil
}

// CreateKey - create a keyvault key
func CreateKey(kv string, key string) (keyvault.KeyBundle, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	ks := KeySize
	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)
	params := keyvault.KeyCreateParameters{
		Kty:     KeyType,
		KeySize: &ks,
	}
	kb, err := c.CreateKey(context.Background(), baseurl, key, params)
	if err != nil {
		return keyvault.KeyBundle{}, err
	}
	return kb, nil
}

// GetKey - return a secret object
func GetKey(kv string, kn string, kve string) (keyvault.KeyBundle, error) {

	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	s, err := c.GetKey(context.Background(), baseurl, kn, kve)
	if err != nil {
		return keyvault.KeyBundle{}, err
	}

	return s, nil
}
