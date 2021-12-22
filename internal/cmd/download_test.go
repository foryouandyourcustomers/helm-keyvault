package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseUri_File(t *testing.T) {
	assert := assert.New(t)

	parsed, err := parseUri("keyvault+file:///path/to/file")
	assert.Nil(err, "should be nil")
	assert.IsType(&fileUri{}, parsed)

}

func Test_parseUri_Keyvault(t *testing.T) {
	assert := assert.New(t)

	parsed, err := parseUri("keyvault+secrets:///path/to/file")
	assert.Nil(err, "should be nil")
	assert.IsType(&keyvaultUri{}, parsed)

}
