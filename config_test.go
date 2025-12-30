package protoreg

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const separator = string(filepath.Separator)

func TestInitDefaults(t *testing.T) {
	c := Config{}
	assert.Error(t, c.InitDefaults())

	c.ImportPaths = []string{
		"./tests/proto/commonapis",
		"./tests/proto/serviceapis",
	}
	c.Proto = []string{""}
	assert.Error(t, c.InitDefaults())

	c.Proto = []string{"unknown/v1/message.proto"}
	assert.Error(t, c.InitDefaults())

	c.Proto = []string{"service/v1/service.proto"}
	assert.NoError(t, c.InitDefaults())
	assert.Equal(t, []string{
		"service" + separator + "v1" + separator + "service.proto",
	}, c.Proto)
}
