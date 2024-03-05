package article

import "C"
import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	// 依赖 gorm 忽略零值的特性
	// 可读性差
	// err := dao.db.WithContext(ctx).updates(&art)
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	// res.RowsAffected 更新行数
	if res.RowsAffected == 0 {
		// 补充点日志
		return fmt.Errorf("更新失败，可能是创作者非法 id %d ,author_id %d",
			art.Id, art.AuthorId)
	}
	return res.Error
}

func (dao *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 先操作制作库（表），再操作线上库（表）
	// 在事务内，采取闭包形态
	// begin, rollback, commit都不需要操心
	var (
		id = art.Id
	)
	// tx, trx => transaction
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewArticleDAO(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		// 操作线上库（表）
		return txDAO.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
}

func (dao *GORMArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	// 插入
	// OnConflict 指数据冲突
	err := dao.db.Clauses(clause.OnConflict{
		// mysql只关心这个
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
		}),
	}).Create(&art).Error
	//最终语句 Insert xxx On duplicate Key update xxx
	return err
}

func NewArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

// Article 制作库
type Article struct {
	Id int64 `gorm:"primaryKey, autoIncrement"`
	//	长度1024
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=TEXT"`
	// 怎么设计索引？
	// where
	// SELECT * FROM articles WHERE author_id = ? order by `utime` desc;
	// 单独查询某一个 select * from articles where id = ?
	AuthorId int64 `gorm:"index=aid_utime"`
	Ctime    int64
	Utime    int64 `gorm:"index=aid_utime"`
}
