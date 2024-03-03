package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
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
