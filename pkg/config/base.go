package config

type BaseConfig struct {
	AppID string `validate:"required" mapstructure:"app_id" yaml:"app_id"`               // application name
	Env   string `validate:"oneof=development production" mapstructure:"env" yaml:"env"` // runtime environment
}
