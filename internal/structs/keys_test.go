package structs

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewKey(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}

	key := NewKey(mock, "mykey", "myversion")

	assert.Equal(KeyvaultObjectId(fmt.Sprintf("https://mykeyvault.%s/keys/mykey/myversion", azure.PublicCloud.KeyVaultDNSSuffix)), key.Kid, "should be equal")
	assert.Equal("mykey", key.Name, "should be equal")
	assert.Equal("myversion", key.Version, "should be equal")
}

func TestKey_Backup(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	key := NewKey(mock, "mykey", "myversion")
	tmpfile, _ := ioutil.TempFile("", "testkey_backup")
	defer os.Remove(tmpfile.Name())
	_ = tmpfile.Close()

	// write backup data (mock keyvault returns name of key as backup content)
	err := key.Backup(tmpfile.Name())
	backup, _ := os.ReadFile(tmpfile.Name())
	assert.Nil(err, "should be nil")
	assert.Equal("mykey", string(backup), "should be equal")
}

func TestKey_Get(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	key := NewKey(mock, "mykey", "myversion")

	key, err := key.Get()

	assert.Nil(err, "should be nil")
	assert.Equal("mykey", key.Name, "should be equal")
	assert.Equal("myversion", key.Version, "should be equal")
	assert.Equal(KeyvaultObjectId(fmt.Sprintf("https://mykeyvault.%s/keys/mykey/myversion", azure.PublicCloud.KeyVaultDNSSuffix)), key.Kid, "should be equal")

}

func TestKey_Create(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	key := NewKey(mock, "mykey", "")
	key, err := key.Create()

	assert.Error(err, "should be errored - mock keyvault doesnt handle multiple keys")
}

func TestKeyList_List(t *testing.T) {
	assert := assert.New(t)

	mock := MockKeyvault{Name: "mykeyvault"}
	kl := KeyList{}
	kl.Keys, _ = kl.List(&mock)

	assert.Len(kl.Keys, 5, "should be 5")
}
