package db

import (
	"fmt"

	gobit "github.com/pot-code/gobit/pkg"
	"go.uber.org/zap"
)

// DBConfig TODO
type DBConfig struct {
	Driver   string   // driver name
	Host     string   // server host
	MaxConn  int32    // maximum opening connections number
	Password string   // db password
	Port     int      // server port
	Protocol string   // connection protocol, eg.tcp
	Query    []string // DSN query parameter
	Schema   string   // use schema
	User     string   // username
	Debug    bool
}

// CreateSqlDBConnection create a DB connection from given config
func CreateSqlDBConnection(cfg *DBConfig, logger *zap.Logger) (conn SqlDB, err error) {
	driver := cfg.Driver
	switch driver {
	case gobit.DriverMysqlDB:
		conn, err = NewMySQLConn(cfg, logger)
	case gobit.DriverPostgresSQL:
		conn, err = NewPostgreSQLConn(cfg, logger)
	default:
		err = fmt.Errorf("unsupported driver: %s", driver)
	}
	return
}

type SqlxDBConfig struct {
	DSN     string
	Driver  string
	MaxConn int32 // maximum opening connections number
	Debug   bool
}

// CreateSqlxDB create a sqlx instance
func CreateSqlxDB(cfg *DBConfig, logger *zap.Logger) (SqlxInterface, error) {
	xconfig := &SqlxDBConfig{Driver: cfg.Driver, MaxConn: cfg.MaxConn, Debug: cfg.Debug}
	if cfg.Driver == gobit.DriverPostgresSQL {
		xconfig.Driver = "pgx"
	}
	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, err
	}
	xconfig.DSN = dsn
	return NewSqlxDB(xconfig, logger)
}
