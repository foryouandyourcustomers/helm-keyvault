package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// TestCreateAndGetKey - create a key via cli command, retrieve it and make sure all values are set correcltu
// try to create a duplicate key which should fail
func (suite *IntegrationTestSuite) TestCreateAndGetKey() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestCreateAndGetKey"

	// test valies
	key := "TestCreateAndGetKey"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")

	// parse the given output, it should be valid json
	createKey, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	suite.Contains(createKey, "kid", "should have 'kid'")
	suite.Contains(createKey, "name", "should have 'name'")
	suite.Contains(createKey, "keyvault", "should have 'keyvault'")
	suite.Contains(createKey, "version", "should have 'version'")
	suite.Equal(key, createKey["name"], "should be equal")
	suite.Equal(suite.AzureKeyVaultName, createKey["keyvault"], "should be equal")
	suite.NotEmpty(createKey["version"], "should not be empty")
	suite.Equal(fmt.Sprintf("https://%s.vault.azure.net/keys/%s/%s", suite.AzureKeyVaultName, key, createKey["version"]), createKey["kid"])

	// execute the create command a second time. this should fail as we only allow creation of a key once
	// to ensure we dont rotate the key by accident (this either needs to be done via a "rotate" argument or via the az cli)
	log.Info("Create duplicate key")
	_, err = runCli(createArgs)
	suite.NotNil(err, "should not be nil")
	suite.Equal(("Key already exists."), error.Error(err))

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	log.Warningln(err)
}

// TestBackupAndRestoreKey - Create a key and a backup file. Delete and restore the key
func (suite *IntegrationTestSuite) TestBackupAndRestoreKey() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestBackupAndRestoreKey"
	// helm-keyvault keys backup --keyvault <keyvaultname> --key "TestBackupAndRestoreKey"

	// test valies
	key := "TestBackupAndRestoreKey"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)
	backupArgs := os.Args[0:1]
	backupArgs = append(backupArgs, "keys", "backup", "--keyvault", suite.AzureKeyVaultName, "--key", key)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")

	// parse the given output, it should be valid json
	createKey, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")

	// create a backup file of the key
	log.Info("Create backup file")
	output, err = runCli(backupArgs)
	suite.Nil(err, "should be nil")

	// check if backup file was created
	fn := fmt.Sprintf("%s.pem", strings.ToUpper(key))
	suite.FileExists(fn, "should exist")

	// delete the key and restore it from the backup file
	// the restore operations is available as function in the keyvault package but
	// not added as a cli operation (yet?)
	log.Info("Remove and restore key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	restoreKey, err := suite.KeyVaultClient.RestoreKey(fn)

	suite.Nil(err, "should be nil")
	suite.Equal(*restoreKey.Key.Kid, createKey["kid"].(string), "should be equal")

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	if err != nil {
		log.Warningln(err)
	}

	// delete backup file
	err = os.Remove(fn)
	if err != nil {
		log.Warningln(err)
	}
}

// TestListKeys - Create a key and list all available keys. As operations run in parallel the test
// succeeds if at least 1 key is retrieved from the keyvault
func (suite *IntegrationTestSuite) TestListKeys() {

	// test cli
	// helm-keyvault keys create --keyvault <keyvaultname> --key "TestListKeys"
	// helm-keyvault keys list --keyvault <keyvaultname>

	// test valies
	key := "TestListKeys"
	createArgs := os.Args[0:1]
	createArgs = append(createArgs, "keys", "create", "--keyvault", suite.AzureKeyVaultName, "--key", key)
	listArgs := os.Args[0:1]
	listArgs = append(listArgs, "keys", "list", "--keyvault", suite.AzureKeyVaultName)

	// execute the create command the first time, this should work ;-)
	log.Info("Create new key")
	output, err := runCli(createArgs)
	suite.Nil(err, "should be nil")

	// list available keys
	log.Info("List available keys")
	output, err = runCli(listArgs)
	suite.Nil(err, "should be nil")

	// parse the given output, it should be valid json
	listKey, err := parseCliOutput(output)
	suite.Nil(err, "should be nil")
	suite.Contains(listKey, "keys", "should have 'keys'")
	suite.GreaterOrEqual(len(listKey["keys"].([]interface{})), 1, "should be greater or equal")

	// delete key
	log.Info("Removing key")
	_, err = suite.KeyVaultClient.Client.DeleteKey(context.Background(), suite.KeyVaultClient.BaseUrl, key)
	if err != nil {
		log.Warningln(err)
	}
}
