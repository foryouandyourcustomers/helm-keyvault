package cmd

import (
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
