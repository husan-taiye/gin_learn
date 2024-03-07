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

	Set(ctx context.Context, id int64, art domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.firstPageKey(uid)).Bytes()
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
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(uid), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleCache) Set(ctx context.Context, id int64, art domain.Article) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	// 过期时间短一些
	return r.client.Set(ctx, r.key(id), data, time.Second*30).Err()
}

func (r *RedisArticleCache) key(id int64) string {
	return fmt.Sprintf("article:id:%d", id)

}

func (r *RedisArticleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}
