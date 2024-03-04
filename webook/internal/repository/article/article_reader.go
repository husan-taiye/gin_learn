package article

import (
	"context"
	"gin_learn/webook/internal/domain"
)

type ArticleReaderRepository interface {
	// 有就更新 没有就新建， upsert
	Save(ctx context.Context, article domain.Article) (int64, error)
}
