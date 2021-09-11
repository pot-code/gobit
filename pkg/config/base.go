package config

type BaseConfig struct {
	AppID string `validate:"required" mapstructure:"app_id" yaml:"app_id"`               // application name
	Port  string `mapstructure:"port" yaml:"port"`                                       // bind listen port
	Env   string `validate:"oneof=development production" mapstructure:"env" yaml:"env"` // runtime environment
}
