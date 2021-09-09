package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddConfigPath(t *testing.T) {
	defaultLength := 1
	cm := NewConfigManager(AddConfigPath("test"))

	assert.Contains(t, cm.configPaths, "test")
	assert.Len(t, cm.configPaths, defaultLength+1)

	cm = NewConfigManager(AddConfigPath("test"), AddConfigPath("test"), AddConfigPath("test"))
	assert.Len(t, cm.configPaths, defaultLength+1)
}

func TestWithConfigName(t *testing.T) {
	cm := NewConfigManager(WithConfigName("test.env"))
	assert.Equal(t, cm.configName, "test.env")
}

func TestWithConfigType(t *testing.T) {
	cm := NewConfigManager(WithConfigType("env"))
	assert.Equal(t, cm.configType, "env")
}

func TestWithEnvPrefix(t *testing.T) {
	cm := NewConfigManager(WithEnvPrefix("GO"))
	assert.Equal(t, cm.envPrefix, "GO")

	cm = NewConfigManager(WithEnvPrefix("TEST"))
	assert.Equal(t, cm.envPrefix, "TEST")
}

func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()
	assert.NotNil(t, cm)
}

func TestConfigManager_LoadConfig(t *testing.T) {
	type TestConfig struct {
		AppID string `mapstructure:"app_id"`
		Port  string `mapstructure:"port"`
	}

	os.Setenv("GO_APP_ID", "test")
	os.Setenv("GO_PORT", ":8081")

	cm := NewConfigManager()
	cfg := new(TestConfig)
	assert.NoError(t, cm.LoadConfig(cfg))

	// see https://github.com/spf13/viper/issues/761
	// assert.Equal(t, "test", cfg.AppID)
	// assert.Equal(t, ":8081", cfg.Port)
}
