package repository

import (
	"context"
	"fmt"
	"lifelog-grpc/user/domain"
	"lifelog-grpc/user/repository/cache"
	"lifelog-grpc/user/repository/dao"
)

var (
	ErrEmailExist   = dao.ErrEmailExist
	ErrUserNotExist = dao.ErrUserNotExist
)

type UserRepository interface {
	CreateUser(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	FindById(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	Update(ctx context.Context, userDomain domain.UserDomain) error
	DeleteByIds(ctx context.Context, ids []int64) error
	FindByEmail(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	Logout(ctx context.Context, sessionId string) error
	SetHistoryPassword(ctx context.Context, email, historyPassword string) error
	GetHistoryPassword(ctx context.Context, email string) ([]string, error)
	LoginOrRegister(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
}

type UserRepositoryV1 struct {
	userDao   dao.UserDao
	userCache cache.UserCache
}

func NewUserRepository(userDao dao.UserDao, userCache cache.UserCache) UserRepository {
	return &UserRepositoryV1{
		userDao:   userDao,
		userCache: userCache,
	}
}

func (u *UserRepositoryV1) LoginOrRegister(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userDao.GetUserByPhoneOrInsert(ctx, userDomain)
}

// GetHistoryPassword 获取历史密码
func (u *UserRepositoryV1) GetHistoryPassword(ctx context.Context, email string) ([]string, error) {
	// 获取用户的 Redis Key（假设用户唯一标识）
	userKey := userKey(email)
	return u.userCache.GetHistoryPassword(ctx, userKey)
}

// SetHistoryPassword 设置历史密码
func (u *UserRepositoryV1) SetHistoryPassword(ctx context.Context, email, historyPassword string) error {
	// 获取用户的 Redis Key（假设用户唯一标识）
	userKey := userKey(email)
	return u.userCache.SetHistoryPassword(ctx, userKey, historyPassword)
}

func userKey(email string) string {
	return fmt.Sprintf("history_passwords:%s", email)
}

func (u *UserRepositoryV1) Logout(ctx context.Context, sessionId string) error {
	return u.userCache.Set(ctx, sessionId)
}

func (u *UserRepositoryV1) CreateUser(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userDao.Insert(ctx, userDomain)
}

func (u *UserRepositoryV1) FindById(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userDao.GetUserById(ctx, userDomain.Id)
}

func (u *UserRepositoryV1) FindByEmail(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userDao.GetUserByEmail(ctx, userDomain.Email)
}

func (u *UserRepositoryV1) Update(ctx context.Context, userDomain domain.UserDomain) error {
	return u.userDao.UpdateById(ctx, userDomain)
}

func (u *UserRepositoryV1) DeleteByIds(ctx context.Context, ids []int64) error {
	return u.userDao.DeleteByIds(ctx, ids)
}
