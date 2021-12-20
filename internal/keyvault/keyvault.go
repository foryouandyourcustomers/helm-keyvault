package keyvault

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"path"
)

const (
	KeyType keyvault.JSONWebKeyType                = "RSA"
	KeySize int32                                  = 4096
	KeyAlgo keyvault.JSONWebKeyEncryptionAlgorithm = keyvault.RSA15
)

// Keyvault Interface - implements Keyvault struct and allows for easier mocking for testing
type KeyvaultInterface interface {
	NewAuthorizer() (autorest.Authorizer, error)
	GetKeyvaultName() string
	// secrets operations
	GetSecret(sn string, sv string) (keyvault.SecretBundle, error)
	PutSecret(name string, value string) (keyvault.SecretBundle, error)
	ListSecrets() ([]keyvault.SecretBundle, error)
	// keys operations
	EncryptString(key string, version string, encoded string) (keyvault.KeyOperationResult, error)
	DecryptString(key string, version string, encrypted string) (keyvault.KeyOperationResult, error)
	ListKeys() ([]keyvault.KeyBundle, error)
	BackupKey(key string) (string, error)
	CreateKey(key string) (keyvault.KeyBundle, error)
	GetKey(key string, version string) (keyvault.KeyBundle, error)
}

type Keyvault struct {
	Client     keyvault.BaseClient
	Authorizer autorest.Authorizer
	Name       string
	BaseUrl    string
}

func (k Keyvault) MarshalJSON() ([]byte, error) {
	name := fmt.Sprintf("\"%s\"", k.Name)
	return []byte(name), nil
}

func (k *Keyvault) GetKeyvaultName() string {
	return k.Name
}

// NewAuthorizer - Returns an authorizer object dependent on config
func (k *Keyvault) NewAuthorizer() (autorest.Authorizer, error) {
	var authorizer autorest.Authorizer
	var err error

	authorizer, err = kvauth.NewAuthorizerFromFile()
	if err != nil {
		authorizer, err = kvauth.NewAuthorizerFromEnvironment()
		if err != nil {
			authorizer, err = kvauth.NewAuthorizerFromCLI()
			if err != nil {
				return nil, errors.New("Unable to initialize keyvault authorizer")
			}
		}
	}

	return authorizer, nil
}

// GetSecret - return a secret object
func (k *Keyvault) GetSecret(sn string, sv string) (keyvault.SecretBundle, error) {

	s, err := k.Client.GetSecret(context.Background(), k.BaseUrl, sn, sv)
	if err != nil {
		return keyvault.SecretBundle{}, err
	}

	return s, nil
}

// PutSecret - put secret into keyvault
func (k *Keyvault) PutSecret(name string, value string) (keyvault.SecretBundle, error) {

	ct := "base64"
	sp := keyvault.SecretSetParameters{
		Value:       &value,
		ContentType: &ct,
	}

	s, err := k.Client.SetSecret(context.Background(), k.BaseUrl, name, sp)
	if err != nil {
		return keyvault.SecretBundle{}, err
	}
	return s, nil
}

// ListSecrets - list all secrets in the specified keyvault
func (k *Keyvault) ListSecrets() ([]keyvault.SecretBundle, error) {

	ctx := context.Background()
	siter, err := k.Client.GetSecretsComplete(ctx, k.BaseUrl, nil)
	if err != nil {
		log.Fatalf("unable to get list of secrets: %v\n", err)
	}

	var s []keyvault.SecretBundle

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.ID)
		b, err := k.Client.GetSecret(context.Background(), k.BaseUrl, key, "")
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
func (k *Keyvault) EncryptString(key string, version string, encoded string) (keyvault.KeyOperationResult, error) {

	ctx := context.Background()
	param := keyvault.KeyOperationsParameters{
		Algorithm: KeyAlgo,
		Value:     &encoded,
	}
	r, err := k.Client.Encrypt(ctx, k.BaseUrl, key, version, param)
	if err != nil {
		return keyvault.KeyOperationResult{}, err
	}

	return r, nil
}

func (k *Keyvault) DecryptString(key string, version string, encrypted string) (keyvault.KeyOperationResult, error) {

	ctx := context.Background()
	param := keyvault.KeyOperationsParameters{
		Algorithm: KeyAlgo,
		Value:     &encrypted,
	}
	r, err := k.Client.Decrypt(ctx, k.BaseUrl, key, version, param)
	if err != nil {
		return keyvault.KeyOperationResult{}, err
	}

	return r, nil
}

// ListKeys - list all keys in the specified keyvault
func (k *Keyvault) ListKeys() ([]keyvault.KeyBundle, error) {

	ctx := context.Background()
	siter, err := k.Client.GetKeysComplete(ctx, k.BaseUrl, nil)
	if err != nil {
		log.Fatalf("unable to get list of keys: %v\n", err)
	}

	var kb []keyvault.KeyBundle

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.Kid)
		b, err := k.Client.GetKey(context.Background(), k.BaseUrl, key, "")
		if err != nil {
			return []keyvault.KeyBundle{}, err
		}

		kb = append(kb, b)
		err = siter.NextWithContext(ctx)
		if err != nil {
			return []keyvault.KeyBundle{}, err
		}
	}

	return kb, nil
}

// BackupKey - Create a backup of a key which can be used for restoring
func (k *Keyvault) BackupKey(key string) (string, error) {

	kb, err := k.Client.BackupKey(context.Background(), k.BaseUrl, key)
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
func (k *Keyvault) CreateKey(key string) (keyvault.KeyBundle, error) {

	ks := KeySize
	params := keyvault.KeyCreateParameters{
		Kty:     KeyType,
		KeySize: &ks,
	}
	kb, err := k.Client.CreateKey(context.Background(), k.BaseUrl, key, params)
	if err != nil {
		return keyvault.KeyBundle{}, err
	}
	return kb, nil
}

// GetKey - return a secret object
func (k *Keyvault) GetKey(key string, version string) (keyvault.KeyBundle, error) {

	s, err := k.Client.GetKey(context.Background(), k.BaseUrl, key, version)
	if err != nil {
		return keyvault.KeyBundle{}, err
	}

	return s, nil
}
