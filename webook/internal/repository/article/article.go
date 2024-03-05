package article

import (
	"context"
	"gin_learn/webook/internal/domain"
	adao "gin_learn/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, art domain.Article) error
}
type CacheArticleRepository struct {
	dao adao.ArticleDAO

	// v1 操作两个dao
	readerDao adao.ReaderDao
	authorDao adao.AuthorDao

	// 耦合了DAO操作的东西
	// 正常条件，想要在repo层面使用事务
	// 那么就只能利用db开启事务之后，创建基于事务的 DAO
	// 或者，直接去掉 DAO 这一层，在repo实现中，直接操作db
	db *gorm.DB
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, art domain.Article) error {
	return c.dao.SyncStatus(ctx, art)
}

func (c *CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, adao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

func (c *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, adao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

// SyncV2 尝试在repository层面解决事务问题
// 确保保存到两个库（表）同时成功，或者同时失败
func (c *CacheArticleRepository) SyncV2(ctx context.Context, article domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	// 利用tx来构建dao
	author := adao.NewAuthorDao(tx)
	reader := adao.NewReaderDao(tx)

	var (
		id  = article.Id
		err error
	)
	artn := c.toEntity(article)
	if id > 0 {
		err = author.UpdateById(ctx, artn)
	} else {
		id, err = author.Insert(ctx, artn)
	}
	if err != nil {
		// 执行有问题，回滚
		// defer以后就不用回滚
		//tx.Rollback()
		return id, err
	}
	// 操作线上库
	// 线上库可能有可能没有 要有一个 upsert的写法
	//err = reader.Upsert(ctx, artn)
	err = reader.UpsertV2(ctx, adao.PublishArticle{
		Article: artn,
	})
	// 执行成功 提交
	tx.Commit()
	return id, err

}

func (c *CacheArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(article))
}

// SyncV1 不同库 不使用事务
func (c *CacheArticleRepository) SyncV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	artn := c.toEntity(article)
	if id > 0 {
		err = c.authorDao.UpdateById(ctx, artn)
	} else {
		id, err = c.authorDao.Insert(ctx, artn)
	}
	if err != nil {
		return id, err
	}
	// 操作线上库
	// 线上库可能有可能没有 要有一个 upsert的写法
	err = c.readerDao.Upsert(ctx, artn)
	return id, err
}

func (c *CacheArticleRepository) toEntity(article domain.Article) adao.Article {
	return adao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	}
}

func NewArticleRepository(dao adao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{
		dao: dao,
	}
}
