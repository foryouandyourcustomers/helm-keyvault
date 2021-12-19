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

// NewSecret - return a Secret Struct
func NewSecret(kv string, secret string, version string) Secret {
	return Secret{
		Id:       NewKeyvaultObjectId(kv, "secrets", secret, version),
		Name:     secret,
		KeyVault: kv,
		Version:  version,
	}
}

// Get - retrieve secret from keyvault
func (s *Secret) Get() (Secret, error) {
	sb, err := keyvault.GetSecret(s.KeyVault, s.Name, s.Version)
	if err != nil {
		return Secret{}, err
	}

	sid := KeyvaultObjectId(*sb.ID)

	return Secret{
		Id:       sid,
		Name:     sid.GetName(),
		KeyVault: sid.GetKeyvault(),
		Version:  sid.GetVersion(),
		Value:    *sb.Value,
	}, nil
}

// Put - put secret into keyvault
func (s *Secret) Put() (Secret, error) {
	sb, err := keyvault.PutSecret(s.KeyVault, s.Name, s.Value)
	if err != nil {
		return Secret{}, err
	}

	sid := KeyvaultObjectId(*sb.ID)

	return Secret{
		Id:       sid,
		Name:     sid.GetName(),
		KeyVault: sid.GetKeyvault(),
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

func (sl *SecretList) List(kv string) ([]Secret, error) {
	sb, err := keyvault.ListSecrets(kv)
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
				KeyVault: soid.GetKeyvault(),
				Version:  soid.GetVersion(),
				// lets not add the secrets value to the list
				//Value:    *s.Value,
			})
	}

	return secrets, nil
}
