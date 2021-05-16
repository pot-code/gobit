package db

import (
	"context"
	"database/sql"
	"time"
)

type TxAccessMode int

type TxDeferrableMode int

// TxOptions Provides a universal option struct across different SQL drivers
type TxOptions struct {
	Isolation      sql.IsolationLevel
	AccessMode     TxAccessMode
	DeferrableMode TxDeferrableMode
}

// SqlRows Provides a universal query result struct across different SQL drivers
type SqlRows interface {
	Next() bool
	Scan(dest ...interface{}) (err error)
	Close() error
}

// SqlDB Universal SQL operation interface, to eliminate the gap between different SQL drivers
type SqlDB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error)
	BeginTx(ctx context.Context, opts *TxOptions) (SqlTx, error)
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}

type SqlTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (SqlRows, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Ping(ctx context.Context) error
}

type SqlxInterface interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Insert(ctx context.Context, query string, args interface{}) (sql.Result, error)
	Close(ctx context.Context) error
}

// CacheDB define a key-value storage interface
type CacheDB interface {
	SetExp(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (bool, error)
	IncrBy(ctx context.Context, key string, val int64) (bool, error)
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

// transaction access mode
const (
	AccessReadOnly TxAccessMode = iota
	AccessReadWrite
)

// transaction defer mode
const (
	Deferrable TxDeferrableMode = iota
	NotDeferrable
)
