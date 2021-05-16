package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	gobit "github.com/pot-code/gobit/pkg"
	"github.com/pot-code/gobit/pkg/logging"
)

func getLogQueryArgs(args []interface{}) []interface{} {
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

func getDSN(cfg *DBConfig) (dsn string, err error) {
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

func sqlTxOptionAdapter(opts *TxOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}
	iso := opts.Isolation
	readOnly := opts.AccessMode == AccessReadOnly
	return &sql.TxOptions{
		Isolation: iso,
		ReadOnly:  readOnly,
	}
}

func pgTxOptionAdapter(opts *TxOptions) pgx.TxOptions {
	if opts == nil {
		return pgx.TxOptions{}
	}
	iso := pgx.TxIsoLevel(strings.ToLower(opts.Isolation.String()))

	var access pgx.TxAccessMode
	if opts.AccessMode == AccessReadOnly {
		access = pgx.ReadOnly
	} else {
		access = pgx.ReadWrite
	}

	var deferrable pgx.TxDeferrableMode
	if opts.DeferrableMode == Deferrable {
		deferrable = pgx.Deferrable
	} else {
		deferrable = pgx.NotDeferrable
	}
	return pgx.TxOptions{
		IsoLevel:       iso,
		AccessMode:     access,
		DeferrableMode: deferrable,
	}
}
