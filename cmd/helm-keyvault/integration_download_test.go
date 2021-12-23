package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// TestDownloadSecret - create a secret and download it for use with the helm downloader plugin
func (suite *IntegrationTestSuite) TestDownloadSecret() {

	// test cli
	// helm-keyvault secret put --keyvault <keyvaultname> --secret "TestDownloadSecret" --file short
	// helm-keyvault download certFile keyFile caFile keyvault+secret://<keyvaultname>.vault.azure.net/secrets/TestDownloadSecret

	// write files with example values
	shortFile, err := ioutil.TempFile(os.TempDir(), "TestDownloadSecret")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(shortFile.Name())
	_, err = shortFile.WriteString(CONTENT_SHORT)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}

	// test values
	secret := "TestDownloadSecret"
	uri := fmt.Sprintf("keyvault+secret://%s.vault.azure.net/secrets/%s", suite.AzureKeyVaultName, secret)
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "secret", "put", "--keyvault", suite.AzureKeyVaultName, "--secret", secret, "--file", shortFile.Name())
	downloadArgs := os.Args[0:1]
	downloadArgs = append(downloadArgs, "download", "certFile", "keyFile", "caFile", uri)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new secret")
	_, err = runCli(createArgs)
	suite.Nil(err, "should be nil")

	// download secret
	log.Info("Download secret")
	output, err := runCli(downloadArgs)
	suite.Equal(string(output), CONTENT_SHORT, "should be equal")

	// delete secret
	log.Info("Removing secret")
	_, err = suite.KeyVaultClient.Client.DeleteSecret(context.Background(), suite.KeyVaultClient.BaseUrl, secret)
	if err != nil {
		log.Warningln(err)
	}
}

// TestDownloadFile - create a key via cli command, encrypt a file and "download" it with the helm downloader plugin
func (suite *IntegrationTestSuite) TestDownloadFile() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestDownloadFile"
	// helm-keyvault files encrypt --keyvault <keyvaultname> --key "TestEncryptAndDecryptFileLong" --file long

	// write files with example values
	longFile, err := ioutil.TempFile(os.TempDir(), "TestDownloadFile")
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
	key := "TestDownloadFile"
	uri := fmt.Sprintf("keyvault+file://%s", longFileEnc)
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)
	encryptArgs := os.Args[0:1]
	encryptArgs = append(encryptArgs, "files", "encrypt", "--keyvault", suite.AzureKeyVaultName, "--key", key, "--file", longFile.Name())
	downloadArgs := os.Args[0:1]
	downloadArgs = append(downloadArgs, "download", "certFile", "keyFile", "caFile", uri)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	_, err = runCli(createArgs)
	suite.Nil(err, "should be nil")

	// encrypt file
	log.Info("Encrypt file")
	_, err = runCli(encryptArgs)
	suite.Nil(err, "should be nil")
	suite.FileExists(longFileEnc, "should exist")
	defer os.Remove(longFileEnc)

	// download file
	log.Info("Download file")
	output, err := runCli(downloadArgs)
	suite.Equal(string(output), CONTENT_LONG, "should be equal")

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	if err != nil {
		log.Warningln(err)
	}
}
