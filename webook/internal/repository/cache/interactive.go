package cache

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/interative_incr_cnt.lua
	luaIncrCnt string
)

type InteractiveCache interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

type RedisInteractiveCache struct {
	client *redis.Cmdable
}
