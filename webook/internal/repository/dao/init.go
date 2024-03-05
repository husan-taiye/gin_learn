package dao

import (
	"gin_learn/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &UserProfile{}, &article.Article{}, &article.PublishArticle{})
}
