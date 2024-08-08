package service

import (
	"context"
	"errors"
	"github.com/ChongYanOvO/go-web-template/internal/domain"
	"github.com/ChongYanOvO/go-web-template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorUserDuplicate             = repository.ErrorUserDuplicate
	ErrorUserNotFound              = repository.ErrorUserNotFound
	ErrorInvalidUserNameOrPassword = errors.New("用户名或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc UserService) Register(ctx context.Context, u *domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc UserService) Login(ctx context.Context, userName string, password string) (*domain.User, error) {
	u, err := svc.repo.FindByUserName(ctx, userName)
	if err != nil {
		if errors.Is(err, ErrorUserNotFound) {
			return u, ErrorUserNotFound
		} else {
			return u, err
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return u, ErrorInvalidUserNameOrPassword
	}
	return u, nil
}

func (svc UserService) Profile(ctx context.Context, userId int64) (*domain.User, error) {
	return svc.repo.FindByUserId(ctx, userId)
}
