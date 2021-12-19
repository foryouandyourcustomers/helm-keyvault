package structs

import (
	"encoding/base64"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
)

//type SecretInterface interface {
//	Get() (Secret, error)
//	Put() (Secret, error)
//	Decode() (string, error)
//}

// NewSecret - return a Secret Struct
func NewSecret(kv keyvault.KeyvaultInterface, secret string, version string) Secret {
	return Secret{
		Id:       NewKeyvaultObjectId(kv.GetKeyvaultName(), "secrets", secret, version),
		Name:     secret,
		KeyVault: kv,
		Version:  version,
	}
}

type Secret struct {
	Id KeyvaultObjectId `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
	//KeyVault string `json:"keyvault,omitempty"`
	KeyVault keyvault.KeyvaultInterface `json:"keyvault,omitempty"`
	Version  string                     `json:"version,omitempty"`
	Value    string                     `json:"value,omitempty"`
}

// Get - retrieve secret from keyvault
func (s *Secret) Get() (Secret, error) {

	sb, err := s.KeyVault.GetSecret(s.Name, s.Version)
	if err != nil {
		return Secret{}, err
	}

	sid := KeyvaultObjectId(*sb.ID)

	return Secret{
		Id:       sid,
		Name:     sid.GetName(),
		KeyVault: s.KeyVault,
		Version:  sid.GetVersion(),
		Value:    *sb.Value,
	}, nil
}

// Put - put secret into keyvault
func (s *Secret) Put() (Secret, error) {

	sb, err := s.KeyVault.PutSecret(s.Name, s.Value)
	if err != nil {
		return Secret{}, err
	}

	sid := KeyvaultObjectId(*sb.ID)

	return Secret{
		Id:       sid,
		Name:     sid.GetName(),
		KeyVault: s.KeyVault,
		Version:  sid.GetVersion(),
		Value:    *sb.Value,
	}, nil
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

func (sl *SecretList) List(kv keyvault.KeyvaultInterface) ([]Secret, error) {

	sb, err := kv.ListSecrets()
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	for _, s := range sb {
		soid := KeyvaultObjectId(*s.ID)
		secrets = append(secrets,
			Secret{
				Id:       soid,
				Name:     soid.GetName(),
				KeyVault: kv,
				Version:  soid.GetVersion(),
				// lets not add the secrets value to the list
				//Value:    *s.Value,
			})
	}

	return secrets, nil
}
