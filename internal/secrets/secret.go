package secrets

import (
	"encoding/base64"
)

type Secret struct {
	Id      string
	Name    string
	Version string `json:",omitempty"`
	Value   string `json:",omitempty"`
}

type SecretList struct {
	Secrets []Secret
}

// Decode - decode the given value from base64 to string
func (s *Secret) Decode() (string, error) {

	dec, err := base64.StdEncoding.DecodeString(s.Value)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
