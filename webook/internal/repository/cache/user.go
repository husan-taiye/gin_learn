package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin_learn/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrUserNotFound = errors.New("用户缓存不存在")

//type CacheV1 interface {
//	Get(ctx context.Context, key string) (any error)
//}

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.UserProfile, error)
	Set(ctx context.Context, user domain.UserProfile) error
}

type RedisUserCache struct {
	// 暂未实现	cache CacheV1
	client     redis.Cmdable
	expiration time.Duration
}

// A用到了B，B一定是接口
// A用到了B，B一定是A的字段
// A用到了B，A绝对不初始化B，而是外面注入

// NewRedisUserCache func NewRedisUserCache(client redis.Cmdable, expiration time.Duration) RedisUserCache {

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.UserProfile, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.UserProfile{}, err
	}
	var u domain.UserProfile
	err = json.Unmarshal(val, &u)
	if u == (domain.UserProfile{}) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, user domain.UserProfile) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.key(user.UserId), val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
