package repository

import (
	"context"
	"github.com/ChongYanOvO/go-web-template/internal/domain"
	"github.com/ChongYanOvO/go-web-template/internal/repository/cache"
	"github.com/ChongYanOvO/go-web-template/internal/repository/dao"
)

var (
	ErrorUserDuplicate = dao.ErrorUserDuplicate
	ErrorUserNotFound  = dao.ErrorUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

func (ur UserRepository) Create(ctx context.Context, u *domain.User) error {
	return ur.dao.Insert(ctx, &dao.User{
		UserName: u.UserName,
		Password: u.Password,
	})
}

func (ur UserRepository) FindByUserName(ctx context.Context, userName string) (*domain.User, error) {
	u, err := ur.dao.FindByUserName(ctx, userName)
	if err != nil {
		return &domain.User{}, err
	}
	return &domain.User{
		ID:       u.ID,
		UserName: u.UserName,
		Password: u.Password,
	}, nil
}

func (ur UserRepository) FindByUserId(ctx context.Context, userId int64) (*domain.User, error) {
	// 先从缓存中获取
	u, err := ur.cache.Get(ctx, userId)
	// 缓存中存在
	if err == nil {
		return u, err
	}
	// 缓存中不存在
	ud, err := ur.dao.FindByUserId(ctx, userId)
	if err != nil {
		return &domain.User{}, err
	}
	user := &domain.User{
		ID:       ud.ID,
		UserName: ud.UserName,
		Password: ud.Password,
	}
	// 异步更新缓存
	go func() {
		err = ur.cache.Set(ctx, user)
		if err != nil {
			return
		}
	}()
	return user, nil
}
