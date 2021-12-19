package structs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitPath(t *testing.T) {
	assert := assert.New(t)

	// valid paths
	validpath1 := "/type/name"
	validpath2 := "/type/name/version"
	// invalid paths
	invalidpath1 := "/type"
	invalidpath2 := "/type/name/version/abcdef"

	var p []string
	var err error

	p, err = splitPath(validpath1)
	assert.Nil(err, "is nil")
	assert.Len(p, 3)
	p, err = splitPath(validpath2)
	assert.Nil(err, "is nil")
	assert.Len(p, 4)

	p, err = splitPath(invalidpath1)
	assert.Nil(p, "is nil")
	assert.Error(err, "is error")
	p, err = splitPath(invalidpath2)
	assert.Nil(p, "is nil")
	assert.Error(err, "is error")
}
