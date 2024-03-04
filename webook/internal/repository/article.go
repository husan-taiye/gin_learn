package repository

import (
	"context"
	"gin_learn/webook/internal/domain"
	"gin_learn/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}
type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func (c *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func (c *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}
