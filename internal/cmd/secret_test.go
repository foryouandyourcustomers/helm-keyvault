package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_GetSecret(t *testing.T) {
	assert := assert.New(t)

	// capture stdout
	// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// execute command
	structs.NewKeyVault = newMockKeyVault
	expectedOutput := "{\"id\":\"https://mykeyvault.vault.azure.net/secrets/yarp/123456\",\"name\":\"yarp\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456\",\"value\":\"Exammple Value\"}"
	err := GetSecret("mykeyvault", "yarp", "123456")
	assert.Nil(err, "should be nil")

	// read in output
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	assert.Equal(expectedOutput, string(out))
}

func Test_PutSecret(t *testing.T) {
	assert := assert.New(t)

	// capture stdout
	// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// create file with clear text value
	value := "Hello, World!"
	valueenc := base64.StdEncoding.EncodeToString([]byte(value))
	tmpfile, _ := ioutil.TempFile("", "Test_PutSecret")
	defer os.Remove(tmpfile.Name())
	_ = tmpfile.Close()

	_ = os.WriteFile(tmpfile.Name(), []byte(value), 644)

	// execute put command, make sure received secret corresponds
	// with the value retrieved from the keyvault
	structs.NewKeyVault = newMockKeyVault
	expectedOutput := fmt.Sprintf("{\"id\":\"https://mykeyvault.vault.azure.net/secrets/yarp/123456\",\"name\":\"yarp\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456\",\"value\":\"%s\"}", valueenc)
	err := PutSecret("mykeyvault", "yarp", tmpfile.Name())
	assert.Nil(err, "should be nil")

	// read in output
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	assert.Equal(expectedOutput, string(out))
}

func Test_ListSecrets(t *testing.T) {
	assert := assert.New(t)

	// capture stdout
	// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// execute put command, make sure received secret corresponds
	// with the value retrieved from the keyvault
	structs.NewKeyVault = newMockKeyVault
	expectedoutput := "{\"secrets\":[{\"id\":\"https://mykeyvault.vault.azure.net/secrets/secret-0/123456789\",\"name\":\"secret-0\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456789\"},{\"id\":\"https://mykeyvault.vault.azure.net/secrets/secret-1/123456789\",\"name\":\"secret-1\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456789\"},{\"id\":\"https://mykeyvault.vault.azure.net/secrets/secret-2/123456789\",\"name\":\"secret-2\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456789\"},{\"id\":\"https://mykeyvault.vault.azure.net/secrets/secret-3/123456789\",\"name\":\"secret-3\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456789\"},{\"id\":\"https://mykeyvault.vault.azure.net/secrets/secret-4/123456789\",\"name\":\"secret-4\",\"keyvault\":{\"Name\":\"mykeyvault\",\"BaseUrl\":\"https://mykeyvault.vault.azure.net\"},\"version\":\"123456789\"}]}"
	err := ListSecrets("mykeyvault")
	assert.Nil(err, "should be nil")

	// read in output
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	assert.Equal(expectedoutput, string(out))
}

func Test_BackupSecret(t *testing.T) {
	assert := assert.New(t)

	// capture stdout
	// https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// setup tmp file for backup content
	tmpfile, _ := ioutil.TempFile("", "Test_BackupSecret")
	defer os.Remove(tmpfile.Name())
	_ = tmpfile.Close()

	// execute put command, make sure received secret corresponds
	// with the value retrieved from the keyvault
	structs.NewKeyVault = newMockKeyVault
	err := BackupSecret("mykeyvault", "yarp", tmpfile.Name())
	assert.Nil(err, "should be nil")

	// read file content
	// mock keyvault returns backup function with secret name as content
	backup, _ := os.ReadFile(tmpfile.Name())

	// read in output
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout

	assert.Equal("", string(out))
	assert.Equal("yarp", string(backup))
}
