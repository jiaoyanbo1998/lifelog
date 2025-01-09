package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lifelog-grpc/user/domain"
	"lifelog-grpc/user/repository/cache"
	"lifelog-grpc/user/repository/dao"
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
	UpdateAvatar(ctx context.Context, userId int64, filePath string) error
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

func (u *UserRepositoryV1) UpdateAvatar(ctx context.Context, userId int64, filePath string) error {
	return u.userDao.UpdateAvatar(ctx, userId, filePath)
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
	userInfo, err := u.userDao.GetUserById(ctx, userDomain.Id)
	if err != nil {
		return domain.UserDomain{}, err
	}
	err = u.userCache.SetUserInfo(ctx, userInfo)
	if err != nil {
		return domain.UserDomain{}, err
	}
	return userInfo, nil
}

func (u *UserRepositoryV1) FindByEmail(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userDao.GetUserByEmail(ctx, userDomain.Email)
}

func (u *UserRepositoryV1) Update(ctx context.Context, userDomain domain.UserDomain) error {
	// 获取redis中存储的用户信息
	userInfo, err := u.userCache.GetUserInfo(ctx, userDomain.Id)
	if err != nil {
		return err
	}
	// 比较值是否相等，只修改不相等的数据，减轻数据库压力
	var user dao.User
	if userDomain.Email != userInfo.Email {
		user.Email = sql.NullString{
			String: userDomain.Email,
			Valid:  userDomain.Email != "",
		}
	}
	if userDomain.Password != userInfo.Password {
		user.Password = userDomain.Password
	}
	if userDomain.Phone != userInfo.Phone {
		user.Phone = sql.NullString{
			String: userDomain.Phone,
			Valid:  userDomain.Phone != "",
		}
	}
	if userDomain.NickName != userInfo.NickName {
		user.NickName = userDomain.NickName
	}
	if userDomain.Authority != userInfo.Authority {
		user.Authority = userDomain.Authority
	}
	return u.userDao.UpdateById(ctx, user)
}

func (u *UserRepositoryV1) DeleteByIds(ctx context.Context, ids []int64) error {
	return u.userDao.DeleteByIds(ctx, ids)
}
