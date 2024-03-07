package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_learn/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, uid int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.key(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(data, &arts)
	return arts, err
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return nil
	}
	return r.client.Set(ctx, r.key(uid), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleCache) key(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)

}
