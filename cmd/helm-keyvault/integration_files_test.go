package main

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

// TestEncryptAndDecryptFileShort - create a key via cli command, encrypt and decrypt a file
func (suite *IntegrationTestSuite) TestEncryptAndDecryptFileShort() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestEncryptAndDecryptFileShort"
	// helm-keyvault files encrypt --keyvault <keyvaultname> --key "TestEncryptAndDecryptFileShort" --file short

	// write files with example values
	shortFile, err := ioutil.TempFile(os.TempDir(), "TestEncryptAndDecryptFileShort")
	shortFileEnc := fmt.Sprintf("%s.enc", shortFile.Name())
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(shortFile.Name())
	_, err = shortFile.WriteString(CONTENT_SHORT)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}

	// test values
	key := "TestEncryptAndDecryptFileShort"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)
	encryptArgs := os.Args[0:1]
	encryptArgs = append(encryptArgs, "files", "encrypt", "--keyvault", suite.AzureKeyVaultName, "--key", key, "--file", shortFile.Name())
	decryptArgs := os.Args[0:1]
	decryptArgs = append(decryptArgs, "files", "decrypt", "--file", shortFileEnc)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")
	createKey, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")

	// encrypt file
	log.Info("Encrypt file")
	_, err = runCli(encryptArgs)
	suite.Nil(err, "should be nil")
	suite.FileExists(shortFileEnc, "should exist")
	defer os.Remove(shortFileEnc)

	// parse encrypted file as json
	// the kid in the file should correspond to the kid of the created key
	// the amount of chunks should be 1
	// the lastmodified field should be a timestamp
	var parsed map[string]interface{}
	fc, _ := os.ReadFile(shortFileEnc)
	err = json.Unmarshal(fc, &parsed)
	timestamp, _ := time.Parse("2006-01-02T15:04:05Z07:0", parsed["lastmodified"].(string))
	suite.Nil(err, "should be nil")
	suite.Equal(createKey["kid"].(string), parsed["kid"].(string), "should be equal")
	suite.Equal(len(parsed["chunks"].([]interface{})), 1, "should be equal")
	suite.IsType(time.Time{}, timestamp)

	// with the encrypted file verified lets decrypt it and make sure its
	// content is the same as the original one
	log.Info("Decrypt file")
	// remove the existing file
	_ = os.Remove(shortFile.Name())
	_, err = runCli(decryptArgs)
	suite.Nil(err, "should be nil")
	suite.FileExists(shortFile.Name(), "should exist")
	dec, err := os.ReadFile(shortFile.Name())
	suite.Equal(string(dec), CONTENT_SHORT, "should be equal")

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	if err != nil {
		log.Warningln(err)
	}
}

// TestEncryptAndDecryptFileShort - create a key via cli command, encrypt and decrypt a file
func (suite *IntegrationTestSuite) TestEncryptAndDecryptFileLong() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestEncryptAndDecryptFileLong"
	// helm-keyvault files encrypt --keyvault <keyvaultname> --key "TestEncryptAndDecryptFileLong" --file long

	// write files with example values
	longFile, err := ioutil.TempFile(os.TempDir(), "TestEncryptAndDecryptFileShort")
	longFileEnc := fmt.Sprintf("%s.enc", longFile.Name())
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(longFile.Name())
	_, err = longFile.WriteString(CONTENT_LONG)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}

	// test values
	key := "TestEncryptAndDecryptFileLong"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)
	encryptArgs := os.Args[0:1]
	encryptArgs = append(encryptArgs, "files", "encrypt", "--keyvault", suite.AzureKeyVaultName, "--key", key, "--file", longFile.Name())
	decryptArgs := os.Args[0:1]
	decryptArgs = append(decryptArgs, "files", "decrypt", "--file", longFileEnc)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")
	createKey, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")

	// encrypt file
	log.Info("Encrypt file")
	_, err = runCli(encryptArgs)
	suite.Nil(err, "should be nil")
	suite.FileExists(longFileEnc, "should exist")
	defer os.Remove(longFileEnc)

	// parse encrypted file as json
	// the kid in the file should correspond to the kid of the created key
	// the amount of chunks should be 1
	// the lastmodified field should be a timestamp
	var parsed map[string]interface{}
	fc, _ := os.ReadFile(longFileEnc)
	err = json.Unmarshal(fc, &parsed)
	timestamp, _ := time.Parse("2006-01-02T15:04:05Z07:0", parsed["lastmodified"].(string))
	suite.Nil(err, "should be nil")
	suite.Equal(createKey["kid"].(string), parsed["kid"].(string), "should be equal")
	suite.Equal(len(parsed["chunks"].([]interface{})), 2, "should be equal")
	suite.IsType(time.Time{}, timestamp)

	// with the encrypted file verified lets decrypt it and make sure its
	// content is the same as the original one
	log.Info("Decrypt file")
	// remove the existing file
	_ = os.Remove(longFile.Name())
	_, err = runCli(decryptArgs)
	suite.Nil(err, "should be nil")
	suite.FileExists(longFile.Name(), "should exist")
	dec, err := os.ReadFile(longFile.Name())
	suite.Equal(string(dec), CONTENT_LONG, "should be equal")

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	if err != nil {
		log.Warningln(err)
	}
}
