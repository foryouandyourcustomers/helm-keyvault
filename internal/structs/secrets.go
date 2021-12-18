package structs

import (
	"encoding/base64"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
)

type Secret struct {
	Id       KeyvaultObjectId `json:"id,omitempty"`
	Name     string           `json:"name,omitempty"`
	KeyVault string           `json:"keyvault,omitempty"`
	Version  string           `json:"version,omitempty"`
	Value    string           `json:"value,omitempty"`
}

// Get - retrieve secret from keyvault
func (s *Secret) Get() (string, error) {
	sb, err := keyvault.GetSecret(s.KeyVault, s.Name, s.Version)
	if err != nil {
		return "", err
	}

	s.Value = *sb.Value
	s.Id = KeyvaultObjectId(*sb.ID)

	return s.Value, nil
}

// Put - put secret into keyvault
func (s *Secret) Put() (string, error) {
	sb, err := keyvault.PutSecret(s.KeyVault, s.Name, s.Value)
	if err != nil {
		return "", err
	}
	s.Id = KeyvaultObjectId(*sb.ID)
	return string(s.Id), nil
}

// Decode - decode the given value from base64 to string
func (s *Secret) Decode() (string, error) {

	dec, err := base64.StdEncoding.DecodeString(s.Value)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}

type SecretList struct {
	Secrets []Secret `json:"secrets,omitempty"`
}

func (sl *SecretList) List(kv string) error {
	sb, err := keyvault.ListSecrets(kv)
	if err != nil {
		return err
	}

	for _, s := range sb {
		soid := KeyvaultObjectId(*s.ID)
		sn, _ := soid.GetName()
		skv, _ := soid.GetKeyvault()
		sve, _ := soid.GetVersion()

		sl.Secrets = append(sl.Secrets,
			Secret{
				Id:       soid,
				Name:     sn,
				KeyVault: skv,
				Version:  sve,
				// lets not add the secrets value to the list
				//Value:    *s.Value,
			})
	}

	return nil
}
