package dao

import (
	"context"
	"database/sql"
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

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
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
func (dao *UserDao) Update(ctx context.Context, up UserProfile) error {
	now := time.Now().UnixMilli()
	up.CreateTime = now
	up.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&up).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			var findUp UserProfile
			dao.db.First(&findUp, "user_id = ?", up.UserId)
			findUp.UpdateTime = now
			findUp.Birthday = up.Birthday
			findUp.Nickname = up.Nickname
			findUp.Profile = up.Profile
			dao.db.Save(&findUp)
			//err := dao.db.WithContext(ctx).Where("user_id = ？", up.UserId).Updates(map[string]interface{}{
			//	"nickname": up.Nickname, "birthday": up.Birthday, "profile": up.Profile}).Error
			return nil
		}
	}
	return err
}

func (dao *UserDao) FindProfileByUserId(ctx context.Context, userId int64) (UserProfile, error) {
	var up UserProfile
	err := dao.db.WithContext(ctx).First(&up, "user_id = ?", userId).Error
	return up, err
}

// User 直接对应数据库表结构
// entity/model/PO（persistent object）
type User struct {
	// 数据库模型
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string

	CreateTime int64
	UpdateTime int64
}

type UserProfile struct {
	Id       int64 `gorm:"primaryKey, autoIncrement"`
	UserId   int64 `gorm:"unique"`
	Nickname string

	Birthday   int64
	Profile    string
	CreateTime int64
	UpdateTime int64
}
