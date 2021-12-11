package cmd

import (
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// Download - downlaod the secret from the given keyvault uri
func Download(uri string) {
	log.Debugf("Download keyvault secret with uri %s", uri)

	// parse the given uri
	// keyvault://<keyvaultname>/secrets/<secretname>/<version>
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	// retrieve keyvault, secret and optional secret version from parsed uri
	kv := u.Host
	s := strings.Split(strings.Replace(u.Path, "/secrets/", "", 1), "/")
	sn := s[0]
	sv := ""
	if len(s) == 2 {
		sv = s[1]
	}

	// get secret from keyvault and print it
	secret := keyvault.GetSecret(kv, sn, sv)
	fmt.Print(secret.Decode())
}
