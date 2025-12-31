package protoreg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDefaults(t *testing.T) {
	c := Config{}
	assert.Error(t, c.InitDefaults())

	c.ProtoPath = []string{
		"./tests/proto/commonapis",
		"./tests/proto/serviceapis",
	}
	c.Files = []string{""}
	assert.Error(t, c.InitDefaults())

	c.Files = []string{"unknown/v1/message.proto"}
	assert.Error(t, c.InitDefaults())

	c.Files = []string{"service/v1/service.proto"}
	assert.NoError(t, c.InitDefaults())
	assert.Equal(t, []string{
		"service/v1/service.proto",
	}, c.Files)
}
