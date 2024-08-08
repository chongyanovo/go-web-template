package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrorUserDuplicate = errors.New("用户名冲突")
	ErrorUserNotFound  = gorm.ErrRecordNotFound
)

type User struct {
	ID          int64  `gorm:"primary_key;auto_increment;not_null" json:"id"`
	UserName    string `gorm:"type:varchar(255);not null;unique" json:"userName"`
	Password    string `gorm:"type:varchar(255);not null" json:"password"`
	CreatedTime int64  `gorm:"type:bigint;not null" json:"createdTime"`
	UpdatedTime int64  `gorm:"type:bigint;not null" json:"updatedTime"`
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud UserDao) Insert(ctx context.Context, u *User) error {
	now := time.Now().UnixMilli()
	u.CreatedTime = now
	u.UpdatedTime = now
	err := ud.db.WithContext(ctx).Create(u).Error
	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		const uniqueConflictsErrorCode uint16 = 1062
		switch mysqlError.Number {
		case uniqueConflictsErrorCode:
			return ErrorUserDuplicate
		default:
			return errors.New("其他错误")
		}
	}
	return ud.db.WithContext(ctx).Create(u).Error
}

func (ud UserDao) FindByUserName(ctx context.Context, userName string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("user_name = ?", userName).First(&u).Error
	return u, err
}

func (ud UserDao) FindByUserId(ctx context.Context, userId int64) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("id = ?", userId).First(&u).Error
	return u, err
}
