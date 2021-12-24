package main

import (
	"context"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

// TestCreateAndGetKey - create a key via cli command, retrieve it and make sure all values are set correcltu
// try to create a duplicate key which should fail
func (suite *IntegrationTestSuite) TestCreateAndGetSecret() {

	// test cli
	// helm-keyvault secret put --keyvault <keyvaultname> --secret "TestCreateAndGetSecret --file short"
	// helm-keyvault secret get --keyvault <keyvaultname> --secret "TestCreateAndGetSecret
	// helm-keyvault secret put --keyvault <keyvaultname> --secret "TestCreateAndGetSecret --file long"

	// write files with example values
	shortFile, err := ioutil.TempFile(os.TempDir(), "TestCreateAndGetSecret")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	longFile, err := ioutil.TempFile(os.TempDir(), "TestCreateAndGetSecret")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(shortFile.Name())
	defer os.Remove(longFile.Name())

	_, err = shortFile.WriteString(CONTENT_SHORT)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}
	_, err = longFile.WriteString(CONTENT_LONG)
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// test values
	secret := "TestCreateAndGetSecret"
	createShortArgs := os.Args[0:1]
	createShortArgs = append(createShortArgs, "secret", "put", "--keyvault", suite.AzureKeyVaultName, "--secret", secret, "--file", shortFile.Name())
	createLongArgs := os.Args[0:1]
	createLongArgs = append(createLongArgs, "secret", "put", "--keyvault", suite.AzureKeyVaultName, "--secret", secret, "--file", longFile.Name())
	getArgs := os.Args[0:1]
	getArgs = append(getArgs, "secret", "get", "--keyvault", suite.AzureKeyVaultName, "--secret", secret)

	// execute the create command the first time, this should work ;-)
	log.Info("Create secret with short file content")
	output, err := runCli(createShortArgs)
	suite.Nil(err, "should be nil")
	// parse the given output, it should be valid json
	createSecretShort, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")

	suite.Contains(createSecretShort, "id", "should have 'id'")
	suite.Contains(createSecretShort, "name", "should have 'name'")
	suite.Contains(createSecretShort, "keyvault", "should have 'keyvault'")
	suite.Contains(createSecretShort, "version", "should have 'version'")
	suite.Contains(createSecretShort, "value", "should have 'value'")
	suite.Equal(secret, createSecretShort["name"], "should be equal")
	suite.Equal(suite.AzureKeyVaultName, createSecretShort["keyvault"], "should be equal")
	suite.NotEmpty(createSecretShort["version"], "should not be empty")
	suite.Equal(fmt.Sprintf("https://%s.vault.azure.net/secrets/%s/%s", suite.AzureKeyVaultName, secret, createSecretShort["version"]), createSecretShort["id"])

	// get the created secret from the keyvault and verify the decoded value
	log.Info("Retrieve secret with short file content and compare")
	output, err = runCli(getArgs)
	suite.Nil(err, "should be nil")
	// parse the given output, it should be valid json
	getSecretShort, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	dec, err := base64.StdEncoding.DecodeString(getSecretShort["value"].(string))
	suite.Nil(err, "should be nil")
	suite.Equal(getSecretShort["value"], createSecretShort["value"], "should be equal")
	suite.Equal(string(dec), CONTENT_SHORT, "should be equal")

	// create the same secret, this should create a new version of the secret
	// with the long string content
	log.Info("Create secret with short file content")
	output, err = runCli(createLongArgs)
	suite.Nil(err, "should be nil")
	// parse the given output, it should be valid json
	createSecretLong, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	suite.NotEqual(createSecretShort["version"], createSecretLong["version"], "should not be equal")
	suite.NotEqual(createSecretShort["value"], createSecretLong["value"], "should not be equal")

	log.Info("Retrieve secret with long file content and compare")
	output, err = runCli(getArgs)
	suite.Nil(err, "should be nil")
	// parse the given output, it should be valid json
	getSecretLong, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	dec, err = base64.StdEncoding.DecodeString(getSecretLong["value"].(string))
	suite.Nil(err, "should be nil")
	suite.Equal(getSecretLong["value"], createSecretLong["value"], "should be equal")
	suite.Equal(string(dec), CONTENT_LONG, "should be equal")

	// delete secret
	log.Info("Removing secret")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, secret)
	if err != nil {
		log.Warningln(err)
	}
}

// TestBackupAndRestoreSecret - Create a secret and a backup file. Delete and restore the secret
func (suite *IntegrationTestSuite) TestBackupAndRestoreSecret() {

	// test cli
	// helm-keyvault secret put --keyvault <keyvaultname> --secret "TestBackupAndRestoreSecret -f short"
	// helm-keyvault secret backup --keyvault <keyvaultname> --secret "TestBackupAndRestoreSecret"

	// write files with example values
	shortFile, err := ioutil.TempFile(os.TempDir(), "TestBackupAndRestoreSecret")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(shortFile.Name())
	_, err = shortFile.WriteString(CONTENT_SHORT)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}

	// test values
	secret := "TestBackupAndRestoreSecret"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "secret", "put", "--keyvault", suite.AzureKeyVaultName, "--secret", secret, "--file", shortFile.Name())
	backupArgs := os.Args[0:1]
	backupArgs = append(backupArgs, "secret", "backup", "--keyvault", suite.AzureKeyVaultName, "--secret", secret)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new secret")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")

	// parse the given output, it should be valid json
	createSecret, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")

	// create a backup file of the key
	log.Info("Create backup file")
	output, err = runCli(backupArgs)
	suite.Nil(err, "should be nil")

	// check if backup file was created
	fn := fmt.Sprintf("%s.pem", strings.ToUpper(secret))
	suite.FileExists(fn, "should exist")

	// delete the key and restore it from the backup file
	// the restore operations is available as function in the keyvault package but
	// not added as a cli operation (yet?)
	log.Info("Remove and restore key")
	_, err = suite.KeyVaultClient.Client.DeleteSecret(context.Background(), suite.KeyVaultClient.BaseUrl, secret)
	restoreSecret, err := suite.KeyVaultClient.RestoreSecret(fn)
	suite.Nil(err, "should be nil")
	suite.Equal(*restoreSecret.ID, createSecret["id"].(string), "should be equal")

	// delete key
	log.Info("Removing secret")
	_, err = suite.KeyVaultClient.Client.DeleteSecret(context.Background(), suite.KeyVaultClient.BaseUrl, secret)
	if err != nil {
		log.Warningln(err)
	}

	// delete backup file
	err = os.Remove(fn)
	if err != nil {
		log.Warningln(err)
	}
}

// TestListSecrets - Create a secret and list all available secrets. As operations run in parallel the test
// succeeds if at least 1 secret is retrieved from the keyvault
func (suite *IntegrationTestSuite) TestListSecrets() {

	// test cli
	// helm-keyvault secret create --keyvault <keyvaultname> --key "TestListSecrets" --file short
	// helm-keyvault secret list --keyvault <keyvaultname>

	// write files with example values
	shortFile, err := ioutil.TempFile(os.TempDir(), "TestListSecrets")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(shortFile.Name())
	_, err = shortFile.WriteString(CONTENT_SHORT)
	if err != nil {
		log.Fatal("Unable to write to file", err)
	}

	// test values
	secret := "TestListSecrets"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "secret", "put", "--keyvault", suite.AzureKeyVaultName, "--secret", secret, "--file", shortFile.Name())
	listArgs := os.Args[0:1]
	listArgs = append(listArgs, "secret", "list", "--keyvault", suite.AzureKeyVaultName)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new secret")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")

	// list available secrets
	log.Info("List available secrets")
	output, err = runCli(listArgs)
	suite.Nil(err, "should be nil")

	// parse the given output, it should be valid json
	listSecret, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	suite.Contains(listSecret, "secrets", "should have 'secrets'")
	suite.GreaterOrEqual(len(listSecret["secrets"].([]interface{})), 1, "should be greater or equal")

	// delete secret
	log.Info("Removing secret")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, secret)
	if err != nil {
		log.Warningln(err)
	}
}
