package secrets

import (
	"encoding/base64"
)

type Secret struct {
	Id      string
	Name    string
	Version string
	Value   string
}

// Decode - decode the given value from base64 to string
func (s *Secret) Decode() string {

	dec, err := base64.StdEncoding.DecodeString(s.Value)
	if err != nil {
		panic(err)
	}

	return string(dec)
}
