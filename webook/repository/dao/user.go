package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱重复")
	ErrUserNotFount       = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// User 直接对应数据库表结构
// entity/model/PO（persistent object）
type User struct {
	// 数据库模型
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	CreateTime int64
	UpdateTime int64
}
