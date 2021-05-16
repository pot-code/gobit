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
