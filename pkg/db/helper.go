package db

import (
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v4"
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
