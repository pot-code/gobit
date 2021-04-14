package db

import (
	"context"
	"database/sql"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	gobit "github.com/pot-code/gobit/pkg"
	"go.uber.org/zap"
)

// mysqlConn Wraps a *sql.db object and provides the implementation of ITransactionalDB.
//
// it uses zap for default logging
type mysqlConn struct {
	db     *sql.DB
	logger *zap.Logger
}

// mysqlTx transaction wrapper
type mysqlTx struct {
	tx     *sql.Tx
	logger *zap.Logger
}

// assertion
var _ TransactionalDB = &mysqlConn{}
var _ TransactionalDB = &mysqlTx{}

// NewMySQLConn Returns a MySQL connection pool
func NewMySQLConn(cfg *DBConfig, logger *zap.Logger) (TransactionalDB, error) {
	dsn, _ := GetDSN(cfg)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(int(cfg.MaxConn))
	return &mysqlConn{conn, logger}, err
}

// BeginTx start a new transaction context
func (mw *mysqlConn) BeginTx(ctx context.Context, opts *TxOptions) (TransactionalDB, error) {
	logger := mw.logger
	startTime := time.Now()

	txConfig := mysqlTxOptionAdapter(opts)
	tx, err := mw.db.BeginTx(ctx, txConfig)
	endTime := time.Now()
	logger.Debug("", zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "BeginTx"),
	)
	if err != nil {
		return nil, &SqlDBError{Err: err}
	}
	return &mysqlTx{tx, logger}, err
}

func (mw *mysqlConn) Commit(ctx context.Context) error {
	return nil
}

func (mw *mysqlConn) Rollback(ctx context.Context) error {
	return nil
}

func (mw *mysqlConn) Ping(ctx context.Context) error {
	return mw.db.PingContext(ctx)
}

func (mw *mysqlConn) Close(ctx context.Context) error {
	return mw.db.Close()
}

func (mw *mysqlConn) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := mw.logger
	startTime := time.Now()

	res, err := mw.db.ExecContext(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Exec"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return res, err
}

func (mw *mysqlConn) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := mw.logger
	startTime := time.Now()

	rows, err := mw.db.QueryContext(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Query"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return rows, err
}

func (mt *mysqlTx) BeginTx(ctx context.Context, opts *TxOptions) (TransactionalDB, error) {
	return nil, gobit.ErrReopenTransaction
}

func (mt *mysqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := mt.logger
	startTime := time.Now()

	res, err := mt.tx.ExecContext(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Exec"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return res, err
}

func (mt *mysqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := mt.logger
	startTime := time.Now()

	rows, err := mt.tx.QueryContext(ctx, query, args...)
	endTime := time.Now()
	logger.Debug(query,
		zap.Duration("duration", endTime.Sub(startTime)),
		zap.String("method", "Query"),
		zap.Any("args", GetLogQueryArgs(args)))
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: GetLogQueryArgs(args)}
	}
	return rows, err
}

func (mt *mysqlTx) Commit(ctx context.Context) error {
	logger := mt.logger
	startTime := time.Now()
	err := mt.tx.Commit()
	endTime := time.Now()
	logger.Debug("Commit", zap.Duration("duration", endTime.Sub(startTime)))
	if err != nil {
		return &SqlDBError{Err: err}
	}
	return err
}

func (mt *mysqlTx) Rollback(ctx context.Context) error {
	logger := mt.logger
	startTime := time.Now()
	err := mt.tx.Rollback()
	endTime := time.Now()
	logger.Debug("RollBack", zap.Duration("duration", endTime.Sub(startTime)))
	if err != nil {
		return &SqlDBError{Err: err}
	}
	return err
}

func (mt *mysqlTx) Ping(ctx context.Context) error {
	return nil
}

func (mt *mysqlTx) Close(ctx context.Context) error {
	return nil
}

func mysqlTxOptionAdapter(opts *TxOptions) *sql.TxOptions {
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
