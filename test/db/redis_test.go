package db

import (
	"context"
	"testing"
	"time"

	"github.com/pot-code/gobit/pkg/db"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestConnection(t *testing.T) {
	client := db.NewRedisClient("localhost", 6379, "")
	assert.Equal(t, client.Ping(ctx), nil)
}

func TestSetNumber(t *testing.T) {
	client := db.NewRedisClient("localhost", 6379, "")
	err := client.Set(ctx, "num", 1)
	assert.Nil(t, err)
	num, err := client.Get(ctx, "num")
	assert.Nil(t, err)
	assert.Equal(t, num, "1")
	ok, err := client.Delete(ctx, "num")
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSetExp(t *testing.T) {
	client := db.NewRedisClient("localhost", 6379, "")
	err := client.SetExp(ctx, "num", 1, 100*time.Millisecond)
	assert.Nil(t, err)

	time.Sleep(200 * time.Millisecond)

	num, err := client.Get(ctx, "num")
	assert.Nil(t, err)
	assert.Empty(t, num)
}
