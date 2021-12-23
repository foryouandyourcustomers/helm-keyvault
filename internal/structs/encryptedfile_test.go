package structs

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEncryptedFile_LoadFile_SingleChunk(t *testing.T) {
	assert := assert.New(t)

	// file contents
	content := "My raw string"
	content_decoded := base64.RawURLEncoding.EncodeToString([]byte(content))

	// write tempfile
	tmpfile, _ := ioutil.TempFile("", "TestEncryptedFile_LoadFile_SingleChunk")
	defer os.Remove(tmpfile.Name())
	_, _ = tmpfile.WriteString(content)
	_ = tmpfile.Close()

	// load file
	encfile := EncryptedFile{}
	encoded, err := encfile.LoadFile(tmpfile.Name())

	assert.Nil(err, "should be nil")
	assert.Len(encoded, 1, "should be 1")
	assert.Equal(content_decoded, encoded[0], "should be equal")
}

func TestEncryptedFile_LoadFile_MultipleChunks(t *testing.T) {
	assert := assert.New(t)

	// generate a string with Nx the chunk size ((4096 - 384) / 8) + 6)
	// https://www.admfactory.com/how-to-generate-a-fixed-length-random-string-using-golang/
	chunklen := 3
	chunksize := ((4096 - 384) / 8) + 6

	var content []string
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < chunklen; i++ {
		b := make([]rune, chunksize)
		for i := range b {
			b[i] = letter[rand.Intn(len(letter))]
		}
		content = append(content, string(b))
	}

	// write tempfile
	tmpfile, _ := ioutil.TempFile("", "TestEncryptedFile_LoadFile_MultipleChunks")
	defer os.Remove(tmpfile.Name())
	_, _ = tmpfile.WriteString(strings.Join(content, ""))
	_ = tmpfile.Close()

	// load file
	encfile := EncryptedFile{}
	encoded, err := encfile.LoadFile(tmpfile.Name())

	assert.Nil(err, "should be nil")
	assert.Len(encoded, chunklen, "should be N")
	for i, s := range content {
		e, err := base64.RawURLEncoding.DecodeString(encoded[i])
		assert.Nil(err, "should be nil")
		assert.Equal(s, string(e), "should be equal")

	}
	//assert.Equal(content_decoded, encoded[0], "should be equal")
}

func TestEncryptedFile_LoadEncryptedFile(t *testing.T) {
	assert := assert.New(t)

	// file contents
	content := `{
 		"kid": "https://mykeyvault.vault.azure.net/keys/mykey/myversion",
 		"chunks": [
  			"chunk1",
			"chunk2",
			"chunk3"
 		],
 		"lastmodified": "2021-12-20T21:11:28+01:0"
	}`

	// write tempfile
	tmpfile, _ := ioutil.TempFile("", "TestEncryptedFile_LoadEncryptedFile")
	defer os.Remove(tmpfile.Name())
	_, _ = tmpfile.WriteString(content)
	_ = tmpfile.Close()

	// load file
	encfile := EncryptedFile{}
	encfile, err := encfile.LoadEncryptedFile(tmpfile.Name())

	assert.Nil(err, "should be nil")
	assert.Len(encfile.EncryptedData, 3)
	assert.IsType(JTime{}, encfile.LastModified)
	assert.IsType(KeyvaultObjectId(""), encfile.Kid)
	assert.Equal("mykeyvault", encfile.Kid.GetKeyvault())
	assert.Equal("myversion", encfile.Kid.GetVersion())
	assert.Equal("mykey", encfile.Kid.GetName())
}

func TestEncryptedFile_WriteFile(t *testing.T) {
	assert := assert.New(t)

	// file content
	encfile := EncryptedFile{
		EncodedData: []string{
			"TXkgU3RyaW5nCg", //"My String\n"
			"TXkgU3RyaW5nCg", //"My String\n"
			"TXkgU3RyaW5nCg", //"My String\n"
		},
	}

	// write tempfile
	tmpfile, _ := ioutil.TempFile("", "TestEncryptedFile_WriteFile")
	defer os.Remove(tmpfile.Name())
	_ = tmpfile.Close()

	// write file with "decoded" values to disk
	err := encfile.WriteFile(tmpfile.Name())

	// load file content
	writencontent, _ := os.ReadFile(tmpfile.Name())

	assert.Nil(err, "should be nil")
	assert.Equal("My String\nMy String\nMy String\n", string(writencontent), "should be equal")
}

func TestEncryptedFile_WriteEncryptedFile(t *testing.T) {
	assert := assert.New(t)

	// file content
	encfile := EncryptedFile{
		Kid: KeyvaultObjectId("https://mykeyvault.vault.azure.net/keys/mykey/myversion"),
		EncryptedData: []string{
			"chunk",
			"chunk",
			"chunk",
		},
		LastModified: JTime(time.Now()),
	}

	// write tempfile
	tmpfile, _ := ioutil.TempFile("", "TestEncryptedFile_WriteEncryptedFile")
	//defer os.Remove(tmpfile.Name())
	_ = tmpfile.Close()

	// write file with "decoded" values to disk
	err := encfile.WriteEncryptedFile(tmpfile.Name())

	// load file content as struct
	encfilewritten := EncryptedFile{}
	encfilewritten, errread := encfilewritten.LoadEncryptedFile(fmt.Sprintf("%s.enc", tmpfile.Name()))

	assert.Nil(err, "should be nil")
	assert.Nil(errread, "should be nil")
	assert.Equal(encfilewritten.EncryptedData[0], encfile.EncryptedData[0])
	assert.Equal(encfilewritten.EncryptedData[1], encfile.EncryptedData[1])
	assert.Equal(encfilewritten.EncryptedData[2], encfile.EncryptedData[2])
}
