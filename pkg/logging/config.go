package logging

type LoggingConfig struct {
	FilePath string `mapstructure:"file_path" yaml:"file_path"`
	Level    string `mapstructure:"level" yaml:"level" validate:"oneof=debug info warn error fatal panic dpanic"`
	Format   string `validate:"oneof=json console" mapstructure:"format" yaml:"format"`
}
