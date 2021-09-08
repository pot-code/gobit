package db

// DB status
const (
	StatusDelete = iota
	StatusValid
)

const (
	DriverMysqlDB  = "mysql"
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite"
	DriverPgx      = "pgx"
)
