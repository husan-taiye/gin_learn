package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User 直接对应数据库表结构
// entity/model/PO（persistent object）
type User struct {
	// 数据库模型
	Id       int64
	Email    string
	Password string

	CreateTime int64
	UpdateTime int64
}
