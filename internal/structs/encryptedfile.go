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

// splitChunk - splits the given string into chunks with the given size
func (e *EncryptedFile) splitChunk(original string, size int) []string {
	// https://stackoverflow.com/questions/35179656/slice-chunking-in-go
	var value []string
	for i := 0; i < len(original); i += size {
		end := i + size
		if end > len(original) {
			end = len(original)
		}

		value = append(value, original[i:end])
	}
	return value
}

// LoadFile - Read and base64 encode the given file
func (e *EncryptedFile) LoadFile(f string) ([]string, error) {
	// load file and encode is as base64 url (non padded)
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// calculate max size (based on 4096bit keys) for data chunks
	// https://stackoverflow.com/questions/1496793/rsa-encryption-getting-bad-length
	chunksize := ((4096 - 384) / 8) + 6

	var value []string
	for _, val := range e.splitChunk(string(c), chunksize) {
		value = append(value, base64.RawURLEncoding.EncodeToString([]byte(val)))
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
func (e *EncryptedFile) EncryptData(kv keyvault.KeyvaultInterface, key string, version string) ([]string, error) {

	// loop trough the chunked encoded data strings
	var value []string
	for _, d := range e.EncodedData {
		enc, err := kv.EncryptString(key, version, d)
		if err != nil {
			return nil, err
		}
		value = append(value, *enc.Result)
	}

	return value, nil
}

func (e *EncryptedFile) DecryptData(kv keyvault.KeyvaultInterface, key string, version string) ([]string, error) {

	// decrypt encrypted data chunks
	var value []string
	for _, chunk := range e.EncryptedData {
		dec, err := kv.DecryptString(key, version, chunk)
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
