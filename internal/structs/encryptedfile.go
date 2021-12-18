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
func (e *EncryptedFile) LoadFile(f string) error {
	// load file and encode is as base64 url (non padded)
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	// calculate max size (based on 2048bit keys) for data chunks
	// https://stackoverflow.com/questions/1496793/rsa-encryption-getting-bad-length
	chunksize := ((4096 - 384) / 8) + 6

	// split the given encoded data string into chunks
	// https://stackoverflow.com/questions/35179656/slice-chunking-in-go
	for i := 0; i < len(c); i += chunksize {
		end := i + chunksize
		if end > len(c) {
			end = len(c)
		}

		e.EncodedData = append(e.EncodedData, base64.RawURLEncoding.EncodeToString(c[i:end]))
	}
	return nil
}

func (e *EncryptedFile) LoadEncryptedFile(f string) error {
	// load file and parse its content
	c, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(c, e)
	if err != nil {
		return err
	}

	return nil
}

// EncryptData - Encrypt encoded data string
func (e *EncryptedFile) EncryptData() error {

	// get keyvault and key info
	kv, _ := e.Kid.GetKeyvault()
	key, _ := e.Kid.GetName()
	kver, _ := e.Kid.GetVersion()

	// loop trough the chunked encoded data strings
	for _, d := range e.EncodedData {
		enc, err := keyvault.EncryptString(kv, key, kver, d)
		if err != nil {
			return err
		}
		e.EncryptedData = append(e.EncryptedData, *enc.Result)
	}

	return nil
}

func (e *EncryptedFile) DecryptData(kv string, k string, v string) error {

	// check if overwrite values are set for keyvault,
	// key and key version
	kvname := kv
	var err error
	if kvname == "" {
		kvname, err = e.Kid.GetKeyvault()
		if err != nil {
			return err
		}
	}
	key := k
	if key == "" {
		key, err = e.Kid.GetName()
		if err != nil {
			return err
		}
	}
	version := v
	if version == "" {
		version, err = e.Kid.GetVersion()
		if err != nil {
			return err
		}
	}

	// decrypt encrypted data chunks
	for _, chunk := range e.EncryptedData {
		dec, err := keyvault.DecryptString(kvname, key, version, chunk)
		if err != nil {
			return err
		}
		e.EncodedData = append(e.EncodedData, *dec.Result)
	}

	return nil
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

	for _, chunk := range e.EncodedData {
		c, err := base64.RawURLEncoding.DecodeString(chunk)
		if err != nil {
			return err
		}
		_, err = fp.WriteString(string(c))
		if err != nil {
			return err
		}

	}

	return nil
}
