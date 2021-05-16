package db

import (
	"fmt"
	"strings"

	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/logging"
)

func GetLogQueryArgs(args []interface{}) []interface{} {
	logArgs := make([]interface{}, 0, len(args))

	for _, a := range args {
		switch v := a.(type) {
		case []byte:
			a = logging.SanitizeBytes(v)
		case string:
			a = logging.SanitizeString(v)
		}
		logArgs = append(logArgs, a)
	}

	return logArgs
}

func GetDSN(cfg *DBConfig) (dsn string, err error) {
	query := ""
	if len(cfg.Query) > 0 {
		query = "?" + strings.Join(cfg.Query, "&")
	}
	switch cfg.Driver {
	case gobit.DriverMysqlDB:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Schema, query)
	case gobit.DriverPostgresSQL:
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Schema, query)
	default:
		err = fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
	return
}
