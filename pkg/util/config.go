package util

import (
	"errors"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigOption interface {
	apply(*ConfigManager)
}

type optionFunc func(*ConfigManager)

func (f optionFunc) apply(cm *ConfigManager) {
	f(cm)
}

func AddConfigPath(p string) ConfigOption {
	return optionFunc(func(cm *ConfigManager) {
		cm.configPaths[p] = struct{}{}
	})
}

func WithConfigName(name string) ConfigOption {
	return optionFunc(func(cm *ConfigManager) {
		cm.configName = name
	})
}

func WithConfigType(t string) ConfigOption {
	return optionFunc(func(cm *ConfigManager) {
		cm.configType = t
	})
}

func WithEnvPrefix(prefix string) ConfigOption {
	return optionFunc(func(cm *ConfigManager) {
		cm.envPrefix = prefix
	})
}

var (
	ErrNoConfigFileFound = errors.New("no config file found")
)

type ConfigManager struct {
	configName  string
	configType  string
	configPaths map[string]struct{}
	envPrefix   string
}

func NewConfigManager(options ...ConfigOption) *ConfigManager {
	cm := &ConfigManager{
		envPrefix:  "GO",
		configName: "config",
		configType: "yml,",
		configPaths: map[string]struct{}{
			".": {},
		},
	}

	for _, option := range options {
		option.apply(cm)
	}

	// env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(cm.envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(cm.configName) // name of config file (without extension)
	for p := range cm.configPaths {
		viper.AddConfigPath(p)
	}
	return cm
}

// LoadConfig load config into cfg
func (cm *ConfigManager) LoadConfig(cfg interface{}) error {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrNoConfigFileFound
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}
