package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SqlxDB struct {
	Client *sqlx.DB
	debug  bool
	logger *zap.Logger
}

func NewSqlxDB(cfg *SqlDBConfig, logger *zap.Logger) (*SqlxDB, error) {
	db, err := sqlx.Open(cfg.Driver, cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(cfg.MaxConn))
	db.SetMaxIdleConns(int(cfg.MaxConn) >> 2)

	return &SqlxDB{db, cfg.Debug, logger.With(zap.String("event.module", "SqlxDB"))}, err
}

func (sd SqlxDB) BeginTx(ctx context.Context, opts *TxOptions) (SqlxTxDB, error) {
	tx, err := sd.Client.BeginTxx(ctx, sqlTxOptionAdapter(opts))
	if err != nil {
		return SqlxTxDB{}, err
	}
	return SqlxTxDB{conn: tx, debug: sd.debug, logger: sd.logger}, nil
}

func (sd SqlxDB) Close(ctx context.Context) error {
	return sd.Client.Close()
}

func (sd SqlxDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	logger := sd.logger

	startTime := time.Now()
	err := sd.Client.SelectContext(ctx, dest, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Select"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return nil
}

func (sd SqlxDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	logger := sd.logger

	startTime := time.Now()
	err := sd.Client.GetContext(ctx, dest, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Get"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return nil
}

func (sd SqlxDB) Insert(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	logger := sd.logger

	startTime := time.Now()
	res, err := sd.Client.NamedExecContext(ctx, query, args)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Insert"),
			zap.Any("args", args),
		)
	}
	if err != nil {
		return res, &SqlDBError{Err: err, Sql: query, Args: []interface{}{""}}
	}
	return res, err
}

func (sd SqlxDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := sd.logger

	startTime := time.Now()
	res, err := sd.Client.ExecContext(ctx, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "ExecContext"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		return res, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return res, err
}

type SqlxTxDB struct {
	conn   *sqlx.Tx
	debug  bool
	logger *zap.Logger
}

func (sd SqlxTxDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	logger := sd.logger

	startTime := time.Now()
	err := sd.conn.SelectContext(ctx, dest, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Select"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return nil
}

func (sd SqlxTxDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	logger := sd.logger

	startTime := time.Now()
	err := sd.conn.GetContext(ctx, dest, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Get"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return nil
}

func (sd SqlxTxDB) Insert(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	logger := sd.logger

	startTime := time.Now()
	res, err := sd.conn.NamedExecContext(ctx, query, args)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "Insert"),
			zap.Any("args", args),
		)
	}
	if err != nil {
		return res, &SqlDBError{Err: err, Sql: query, Args: []interface{}{""}}
	}
	return res, err
}

func (sd SqlxTxDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	logger := sd.logger

	startTime := time.Now()
	res, err := sd.conn.ExecContext(ctx, query, args...)
	endTime := time.Now()
	if sd.debug {
		logger.Debug(query,
			zap.Duration("duration", endTime.Sub(startTime)),
			zap.String("method", "ExecContext"),
			zap.Any("args", getLogQueryArgs(args)),
		)
	}
	if err != nil {
		return res, &SqlDBError{Err: err, Sql: query, Args: getLogQueryArgs(args)}
	}
	return res, err
}

func (sd SqlxTxDB) Commit(ctx context.Context) error {
	return sd.conn.Commit()
}

func (sd SqlxTxDB) Rollback(ctx context.Context) error {
	return sd.conn.Rollback()
}
