package trace

type TraceConfig struct {
	URL string `validate:"required" mapstructure:"url" yaml:"url"` // agent url
}
