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

//func (s *Secret) ToEnv() string {
//	return fmt.Sprintf("%s=%s", s.Key, s.Value)
//}
//
//type SecretList struct {
//	Secrets []Secret
//}
//
//func (sl *SecretList) TransformToEnvironmentNames() []Secret {
//	sn := make([]Secret, 0)
//	for _, s := range sl.Secrets {
//		s.Key = strings.ToLower(strings.ReplaceAll(s.Key, "-", "_"))
//		sn = append(sn, s)
//	}
//	return sn
//}
//
//func (sl *SecretList) FilterSecretsByPrefix(prefix string) []Secret {
//	sn := make([]Secret, 0)
//	for _, s := range sl.Secrets {
//		if strings.HasPrefix(s.Key, prefix) {
//			sn = append(sn, s)
//		} else {
//			log.Debugf("Secret %s doesnt have prefix %s. Dropping it", s.Key, prefix)
//
//		}
//	}
//	return sn
//}
//
//func (sl *SecretList) RemoveSecretsPrefix(prefix string) []Secret {
//	sn := make([]Secret, 0)
//	for _, s := range sl.Secrets {
//		if strings.HasPrefix(s.Key, prefix) {
//			s.Key = strings.TrimLeftFunc(strings.TrimPrefix(s.Key, prefix), func(r rune) bool {
//				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
//			})
//			sn = append(sn, s)
//		} else {
//			log.Debugf("Secret %s doesnt have prefix %s. No need for removal", s.Key, prefix)
//
//		}
//	}
//	return sn
//}
//
//func (sl *SecretList) ReplaceSecretsPrefix(prefix string, replace string) []Secret {
//	sn := make([]Secret, 0)
//	for _, s := range sl.Secrets {
//		if strings.HasPrefix(s.Key, prefix) {
//			s.Key = strings.Replace(s.Key, prefix, replace, 1)
//			sn = append(sn, s)
//		} else {
//			log.Debugf("Secret %s doesnt have prefix %s. No need for replacing string", s.Key, prefix)
//		}
//	}
//	return sn
//}
//
//func (sl *SecretList) PrefixSecretsWithTfVar() []Secret {
//	sn := make([]Secret, 0)
//	for _, s := range sl.Secrets {
//		s.Key = "TF_VAR_" + s.Key
//		sn = append(sn, s)
//	}
//	return sn
//}
