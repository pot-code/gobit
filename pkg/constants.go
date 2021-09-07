package gobit

type AppContextKey string

const (
	EnvProduction = "production"
	EnvDevelop    = "development"
)

const (
	DriverMysqlDB  = "mysql"
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite"
	DriverPgx      = "pgx"
)
