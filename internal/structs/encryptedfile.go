package structs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/keyvault"
	"io/ioutil"
	"os"
)

type EncryptedFile struct {
	Kid           KeyvaultObjectId `json:"kid,omitempty"`
	EncodedData   []string         `json:"-"`
	EncryptedData []string         `json:"chunks,omitempty"`
	LastModified  JTime            `json:"lastmodified,omitempty"`
}

// LoadFile - Read and base64 encode the given file
func (e *EncryptedFile) LoadFile(f string) ([]string, error) {
	// load file and encode is as base64 url (non padded)
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// calculate max size (based on 2048bit keys) for data chunks
	// https://stackoverflow.com/questions/1496793/rsa-encryption-getting-bad-length
	chunksize := ((4096 - 384) / 8) + 6

	// split the given encoded data string into chunks
	// https://stackoverflow.com/questions/35179656/slice-chunking-in-go
	var value []string
	for i := 0; i < len(c); i += chunksize {
		end := i + chunksize
		if end > len(c) {
			end = len(c)
		}

		value = append(value, base64.RawURLEncoding.EncodeToString(c[i:end]))
	}
	return value, nil
}

func (e *EncryptedFile) LoadEncryptedFile(f string) (EncryptedFile, error) {
	// load file and parse its content
	c, err := os.ReadFile(f)
	if err != nil {
		return EncryptedFile{}, err
	}

	var value EncryptedFile
	err = json.Unmarshal(c, &value)
	if err != nil {
		return EncryptedFile{}, err
	}

	return value, err
}

// EncryptData - Encrypt encoded data strings
func (e *EncryptedFile) EncryptData() ([]string, error) {

	// get keyvault and key info
	kv := e.Kid.GetKeyvault()
	key := e.Kid.GetName()
	kver := e.Kid.GetVersion()

	// loop trough the chunked encoded data strings
	var value []string
	for _, d := range e.EncodedData {
		enc, err := keyvault.EncryptString(kv, key, kver, d)
		if err != nil {
			return nil, err
		}
		value = append(value, *enc.Result)
	}

	return value, nil
}

func (e *EncryptedFile) DecryptData(kv string, k string, v string) ([]string, error) {

	// check if overwrite values are set for keyvault,
	// key and key version
	kvname := kv
	var err error
	if kvname == "" {
		kvname = e.Kid.GetKeyvault()
		if err != nil {
			return nil, err
		}
	}
	key := k
	if key == "" {
		key = e.Kid.GetName()
		if err != nil {
			return nil, err
		}
	}
	version := v
	if version == "" {
		version = e.Kid.GetVersion()
		if err != nil {
			return nil, err
		}
	}

	// decrypt encrypted data chunks
	var value []string
	for _, chunk := range e.EncryptedData {
		dec, err := keyvault.DecryptString(kvname, key, version, chunk)
		if err != nil {
			return nil, err
		}
		value = append(value, *dec.Result)
	}

	return value, nil
}

// WriteFile - Write marshalled file to disk
func (e *EncryptedFile) WriteEncryptedFile(f string) error {
	j, err := json.MarshalIndent(e, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s.enc", f), j, 0644)
	return err
}

func (e *EncryptedFile) WriteFile(f string) error {

	fp, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fp.Close()

	val, err := e.GetDecodedString()
	if err != nil {
		return err
	}
	_, err = fp.WriteString(val)
	if err != nil {
		return err
	}

	return nil
}

// GetDecodedString - Returns the base64 decoded string
func (e *EncryptedFile) GetDecodedString() (string, error) {
	var value string
	for _, chunk := range e.EncodedData {
		c, err := base64.RawURLEncoding.DecodeString(chunk)
		if err != nil {
			return "", err
		}

		value = fmt.Sprintf("%s%s", value, c)
	}
	return value, nil
}
