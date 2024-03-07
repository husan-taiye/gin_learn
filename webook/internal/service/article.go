package service

import (
	"context"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/article"
	"github.com/gin-gonic/gin"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx *gin.Context, Uid int64) (domain.Article, error)
	GetPublishedById(ctx *gin.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo article.ArticleRepository

	// v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
}

func (a *articleService) GetPublishedById(ctx *gin.Context, id int64) (domain.Article, error) {
	return a.repo.GetPublishedById(ctx, id)
}

func (a *articleService) GetById(ctx *gin.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	art.Status = domain.ArticleStatusPrivate
	return a.repo.SyncStatus(ctx, art)
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository, reader article.ArticleReaderRepository) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	// 制作库
	//id, err := a.repo.Create(ctx, art)

	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	// 制作库
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 赋值id 确保制作库跟线上库 id相等
	art.Id = id
	// todo 如果失败重试同步线上库
	return a.reader.Save(ctx, art)
}
