package article

import (
	"context"
	"errors"
	"gin_learn/webook/internal/domain"
	"github.com/gin-gonic/gin"
)

var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, art domain.Article) error
	GetByAuthor(ctx context.Context, uid int64, offset, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx *gin.Context, id int64) (Article, error)
}
