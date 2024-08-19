package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverlayConfigWithEnv(t *testing.T) {
	assert := assert.New(t)
	config := Config{}

	assert.Equal("", config.MQURI)

	os.Setenv("knoperator_MQ_URI", "test")

	err := config.OverlayConfigWithEnv(true)
	assert.Nil(err)
	assert.Equal("test", config.MQURI)
}

func TestGetFakeConfig(t *testing.T) {
	fc := GetFakeConfig()
	assert.NotNil(t, fc)
}
