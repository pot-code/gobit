package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PgExecResult struct {
	ct pgconn.CommandTag
}

type PgQueryResult struct {
	rows pgx.Rows
}

var _ SqlRows = &PgQueryResult{}
var _ sql.Result = &PgExecResult{}

func (pr PgExecResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (pr PgExecResult) RowsAffected() (int64, error) {
	return pr.ct.RowsAffected(), nil
}

func (pr PgQueryResult) Next() bool {
	return pr.rows.Next()
}

func (pr PgQueryResult) Scan(dest ...interface{}) (err error) {
	return pr.rows.Scan(dest...)
}

func (pr PgQueryResult) Close() error {
	pr.rows.Close()
	return nil
}

type pgsql struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// assertion
var _ SqlDB = &pgsql{}

// NewPostgreSQLConn Returns a postgreSQL connection pool
func NewPostgreSQLConn(cfg *DBConfig, logger *zap.Logger) (SqlDB, error) {
	dsn, _ := GetDSN(cfg)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = cfg.MaxConn
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	return &pgsql{conn, logger}, err
}

func (pg *pgsql) BeginTx(ctx context.Context, opts *TxOptions) (SqlTx, error) {
	logger := pg.logger
	startTime := time.Now()

	txConfig := pgTxOptionAdapter(opts)
	tx, err := pg.db.BeginTx(ctx, txConfig)
	endTime := time.Now()
	logger.Debug("",
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "BeginTx"),
	)
	if err != nil {
		return nil, &SqlDBError{Err: err}
	}
	return &postgresTx{tx, logger}, err
}

func (pg *pgsql) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

// Close close the whole pool, you better know what you are doing
func (pg *pgsql) Close(ctx context.Context) error {
	pg.db.Close()
	return nil
}

func (pg *pgsql) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := pg.logger
	startTime := time.Now()

	res, err := pg.db.Exec(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Exec"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return &PgExecResult{res}, err
}

func (pg *pgsql) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := pg.logger
	startTime := time.Now()

	rows, err := pg.db.Query(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Query"),
		zap.Any("args", GetLogQueryArgs(args)),
	)
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return &PgQueryResult{rows}, err
}

type postgresTx struct {
	tx     pgx.Tx
	logger *zap.Logger
}

var _ SqlTx = &postgresTx{}

func (pgt *postgresTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := pgt.logger
	startTime := time.Now()

	res, err := pgt.tx.Exec(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Exec"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return &PgExecResult{res}, err
}

func (pgt *postgresTx) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := pgt.logger
	startTime := time.Now()

	rows, err := pgt.tx.Query(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Query"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return &PgQueryResult{rows}, err
}

func (pgt *postgresTx) Commit(ctx context.Context) error {
	logger := pgt.logger
	startTime := time.Now()
	err := pgt.tx.Commit(ctx)
	endTime := time.Now()
	logger.Debug("Commit",
		zap.Duration("duration", endTime.Sub(startTime)),
	)
	if err != nil {
		return &SqlDBError{Err: err}
	}
	return err
}

func (pgt *postgresTx) Rollback(ctx context.Context) error {
	logger := pgt.logger
	startTime := time.Now()
	err := pgt.tx.Rollback(ctx)
	endTime := time.Now()
	logger.Debug("RollBack", zap.Duration("duration", endTime.Sub(startTime)))
	if err != nil {
		return &SqlDBError{Err: err}
	}
	return err
}

func (pgt *postgresTx) Ping(ctx context.Context) error {
	return pgt.tx.Conn().Ping(ctx)
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
