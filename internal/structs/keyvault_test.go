package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewKeyvaultObjectId(t *testing.T) {
	assert := assert.New(t)

	objectid := NewKeyvaultObjectId("mykeyvault", "mytype", "myname", "myversion")

	assert.Equal("mykeyvault", objectid.GetKeyvault(), "should be equal")
	assert.Equal("mytype", objectid.GetType(), "should be equal")
	assert.Equal("myname", objectid.GetName(), "should be equal")
	assert.Equal("myversion", objectid.GetVersion(), "should be equal")
}

func TestNewKeyvaultObjectIdNoVersion(t *testing.T) {
	assert := assert.New(t)

	objectid := NewKeyvaultObjectId("mykeyvault", "mytype", "myname", "")

	assert.Empty(objectid.GetVersion(), "should be empty")
}
