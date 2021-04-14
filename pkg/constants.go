package gobit

type AppContextKey string

const (
	DriverMysqlDB     = "mysql"
	DriverPostgresSQL = "postgres"
)

const (
	EnvProduction = "production"
	EnvDevelop    = "develop"
)

const (
	DefaultLoggingContextKey = AppContextKey("logger")
	DefaultUserContextKey    = AppContextKey("user")
	DefaultLangContextKey    = AppContextKey("lang")
)

const (
	DefaultPaginationEchoKey   = "pagination"
	DefaultRefreshTokenEchoKey = "refresh"
	DefaultPaginationLimit     = 5
	DefaultPaginationOffset    = 0
	SanitizeStringLength       = 64
	SanitizeBytesLength        = 64
	SanitizedSuffix            = "[truncated]"
)
