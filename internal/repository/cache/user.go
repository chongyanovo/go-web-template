package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ChongYanOvO/go-web-template/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrorKeyNotExist = errors.New("数据在缓存中为未实现")

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (uc UserCache) NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (uc UserCache) Set(ctx context.Context, u *domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := uc.generateKey(u.ID)
	return uc.client.Set(ctx, key, val, uc.expiration).Err()
}

func (uc UserCache) Get(ctx context.Context, userId int64) (*domain.User, error) {
	key := uc.generateKey(userId)
	val, err := uc.client.Get(ctx, key).Bytes()
	if err != nil {
		return &domain.User{}, ErrorKeyNotExist
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return &u, err
}

func (uc UserCache) generateKey(userId int64) string {
	return fmt.Sprintf("user:info:%d", userId)
}
