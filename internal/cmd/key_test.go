package cmd

import (
	"github.com/foryouandyourcustomers/helm-keyvault/internal/structs"
	"testing"
)

func Test_GetKey(t *testing.T) {
	//assert := assert.New(t)

	structs.NewKeyVault = newMockKeyVault

	_ = GetSecret("mykeyvault", "yarp", "")
}
