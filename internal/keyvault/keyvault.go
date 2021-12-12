package keyvault

import (
	"context"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
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
func GetSecret(kv string, sn string, sv string) secrets.Secret {

	c := keyvault.New()
	c.Authorizer = authorizer
	s, err := c.GetSecret(context.Background(), fmt.Sprintf("https://%s", kv), sn, sv)
	if err != nil {
		panic(err)
	}

	return secrets.Secret{
		*s.ID,
		sn,
		path.Base(*s.ID),
		*s.Value,
	}
}
