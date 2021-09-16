package db

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pot-code/gobit/pkg/util"
)

func NewSqlxProvider(dc *DatabaseConfig, lm *util.LifecycleManager) *sqlx.DB {
	if dc == nil {
		panic("DatabaseConfig is nil")
	}

	conn, err := sqlx.Connect(dc.Driver, dc.Dsn)
	util.HandlePanicError("failed to create DB connection", err)

	conn.DB.SetMaxOpenConns(int(dc.MaxConn))
	conn.DB.SetMaxIdleConns(int(dc.MaxConn) >> 2)

	lm.AddLivenessProbe(func(ctx context.Context) error {
		return conn.PingContext(ctx)
	})
	lm.OnExit(func(ctx context.Context) {
		log.Println("[sqlx.DB] close DB connection")
		conn.Close()
	})
	return conn
}

func NewRedisCacheProvider(cc *CacheConfig, lm *util.LifecycleManager) *redis.Client {
	if cc == nil {
		panic("CacheConfig is nil")
	}

	addr := fmt.Sprintf("%s:%d", cc.Host, cc.Port)
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cc.Password,
	})

	lm.AddLivenessProbe(func(ctx context.Context) error {
		return rc.Ping(ctx).Err()
	})
	lm.OnExit(func(ctx context.Context) {
		log.Println("[redis.Client] close connection")
		rc.Close()
	})
	return rc
}
