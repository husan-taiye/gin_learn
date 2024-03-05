package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDao interface {
	Upsert(ctx context.Context, article Article) error
	UpsertV2(ctx context.Context, article PublishArticle) error
}

// PublishArticle 线上表
type PublishArticle struct {
	Article
}

func NewReaderDao(db *gorm.DB) ReaderDao {
	panic("")
}
