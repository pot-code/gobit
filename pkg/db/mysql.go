package db

import (
	"context"
	"database/sql"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// mysqlConn Wraps a *sql.db object and provides the implementation of ITransactionalDB.
//
// it uses zap for default logging
type mysqlConn struct {
	db     *sql.DB
	debug  bool
	logger *zap.Logger
}

// assertion
var _ SqlDB = &mysqlConn{}

// NewMySQLConn Returns a MySQL connection pool
func NewMySQLConn(cfg *DBConfig, logger *zap.Logger) (SqlDB, error) {
	dsn, _ := getDSN(cfg)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(int(cfg.MaxConn))
	return &mysqlConn{conn, cfg.Debug, logger}, err
}

// BeginTx start a new transaction context
func (mw *mysqlConn) BeginTx(ctx context.Context, opts *TxOptions) (SqlTx, error) {
	logger := mw.logger
	startTime := time.Now()

	txConfig := sqlTxOptionAdapter(opts)
	tx, err := mw.db.BeginTx(ctx, txConfig)
	endTime := time.Now()
	if mw.debug {
		logger.Debug("", zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "BeginTx"),
		)
	}
	if err != nil {
		return nil, &SqlDBError{Err: err}
	}
	return &mysqlTx{tx, mw.debug, logger}, err
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
	if mw.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Exec"),
			zap.Any("args", getLogQueryArgs(args)))
	}
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return res, err
}

func (mw *mysqlConn) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := mw.logger
	startTime := time.Now()

	rows, err := mw.db.QueryContext(ctx, query, args...)
	endTime := time.Now()
	if mw.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Query"),
			zap.Any("args", getLogQueryArgs(args)))
	}
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return rows, err
}

// mysqlTx transaction wrapper
type mysqlTx struct {
	tx     *sql.Tx
	debug  bool
	logger *zap.Logger
}

var _ SqlTx = &mysqlTx{}

func (mt *mysqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := mt.logger
	startTime := time.Now()

	res, err := mt.tx.ExecContext(ctx, query, args...)
	endTime := time.Now()
	if mt.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Exec"),
			zap.Any("args", getLogQueryArgs(args)))
	}
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return res, err
}

func (mt *mysqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error) {
	logger := mt.logger
	startTime := time.Now()

	rows, err := mt.tx.QueryContext(ctx, query, args...)
	endTime := time.Now()
	if mt.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Query"),
			zap.Any("args", getLogQueryArgs(args)))
	}
	if err != nil {
		return nil, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return rows, err
}

func (mt *mysqlTx) Commit(ctx context.Context) error {
	logger := mt.logger
	startTime := time.Now()
	err := mt.tx.Commit()
	endTime := time.Now()
	if mt.debug {
		logger.Debug("Commit", zap.Duration("duration", endTime.Sub(startTime)))
	}
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
	if mt.debug {
		logger.Debug("RollBack", zap.Duration("duration", endTime.Sub(startTime)))
	}
	if err != nil {
		return &SqlDBError{Err: err}
	}
	return err
}

func (mt *mysqlTx) Ping(ctx context.Context) error {
	return nil
}
