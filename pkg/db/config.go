package db

type DatabaseConfig struct {
	Driver  string `validate:"required" mapstructure:"driver" yaml:"driver"`  // driver name
	Dsn     string `validate:"required" mapstructure:"dsn" yaml:"dsn"`        // dsn string
	MaxConn int32  `validate:"min=1" mapstructure:"max_conn" yaml:"max_conn"` // maximum opening connections number
}

type CacheConfig struct {
	Host     string `validate:"required" mapstructure:"host" yaml:"host"` // bind host address
	Port     int    `validate:"required" mapstructure:"port" yaml:"port"` // bind listen port
	Password string `mapstructure:"password" yaml:"password"`             // password for security reasons
}
