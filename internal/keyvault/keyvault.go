package keyvault

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keys"
	log "github.com/sirupsen/logrus"
	"os"
	"path"

	"github.com/foryouandyourcustomers/helm-keyvault/internal/secrets"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

var (
	authorizer autorest.Authorizer
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
func GetSecret(kv string, sn string, sv string) (secrets.Secret, error) {

	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	s, err := c.GetSecret(context.Background(), baseurl, sn, sv)
	if err != nil {
		return secrets.Secret{}, err
	}

	return secrets.Secret{
		Id:      *s.ID,
		Name:    sn,
		Version: path.Base(*s.ID),
		Value:   *s.Value,
	}, nil
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
func ListSecrets(kv string) (secrets.SecretList, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	ctx := context.Background()
	siter, err := c.GetSecretsComplete(ctx, baseurl, nil)
	if err != nil {
		log.Fatalf("unable to get list of secrets: %v\n", err)
	}

	var s secrets.SecretList

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.ID)
		b, err := c.GetSecret(context.Background(), baseurl, key, "")
		if err != nil {
			return secrets.SecretList{}, err
		}

		s.Secrets = append(s.Secrets, secrets.Secret{Id: *b.ID, Name: path.Base(path.Dir(*b.ID)), Version: path.Base(*b.ID)})
		err = siter.NextWithContext(ctx)
		if err != nil {
			return secrets.SecretList{}, err
		}
	}

	return s, nil
}

// ListKeys - list all keys in the specified keyvault
func ListKeys(kv string) (keys.KeyList, error) {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)

	ctx := context.Background()
	siter, err := c.GetKeysComplete(ctx, baseurl, nil)
	if err != nil {
		log.Fatalf("unable to get list of keys: %v\n", err)
	}

	var k keys.KeyList

	for siter.NotDone() {
		i := siter.Value()

		key := path.Base(*i.Kid)
		b, err := c.GetKey(context.Background(), baseurl, key, "")
		if err != nil {
			return keys.KeyList{}, err
		}

		k.Keys = append(k.Keys, keys.Key{Kid: *b.Key.Kid, Name: path.Base(path.Dir(*b.Key.Kid)), Version: path.Base(*b.Key.Kid)})
		err = siter.NextWithContext(ctx)
		if err != nil {
			return keys.KeyList{}, err
		}
	}

	return k, nil
}

// BackupKey - Create a protected backup file for the import into azure keyvault
func BackupKey(kv string, key string, f string) error {
	c := keyvault.New()
	c.Authorizer = authorizer

	baseurl := fmt.Sprintf("https://%s.%s", kv, azure.PublicCloud.KeyVaultDNSSuffix)
	kb, err := c.BackupKey(context.Background(), baseurl, key)
	if err != nil {
		return err
	}

	fp, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec, err := base64.RawURLEncoding.DecodeString(*kb.Value)
	_, err = fp.WriteString(string(dec))
	if err != nil {
		return err
	}
	return nil
}
