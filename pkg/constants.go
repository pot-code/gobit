package gobit

type AppContextKey string

const (
	EnvProduction = "production"
	EnvDevelop    = "develop"
)

const (
	DriverMysqlDB  = "mysql"
	DriverPostgres = "postgres"
	DriverPgx      = "pgx"
)
