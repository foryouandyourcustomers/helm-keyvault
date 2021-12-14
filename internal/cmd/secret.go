package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// GetSecret - Get secret from given keyvault
func GetSecret(id string) error {

	// parse given id into its keyvault components
	su, err := newSecretUri(id)
	if err != nil {
		return err
	}

	// retrieve and decode secret content
	value, err := su.download()
	if err != nil {
		return err
	}
	fmt.Print(value)
	return nil
}

// PutSecret - Encode file and put secret into keyvault
func PutSecret(id string, f string) error {

	// parse given id into its keyvault components
	su, err := newSecretUri(id)
	if err != nil {
		return err
	}

	// read file and convert it to base64
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	e := base64.StdEncoding.EncodeToString(c)

	// upload the file content'
	value, err := su.upload(e)
	if err != nil {
		return err
	}
	fmt.Print(value)
	return nil
}
