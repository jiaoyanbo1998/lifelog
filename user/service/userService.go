package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"lifelog-grpc/user/domain"
	"lifelog-grpc/user/repository"
	"time"
)

var (
	ErrEmailExist      = repository.ErrEmailExist
	ErrUserNotExist    = repository.ErrUserNotExist
	ErrPasswordWrong   = errors.New("旧密码错误输入错误")
	ErrUseOldPassword  = errors.New("不能使用历史密码")
	NeedUpdatePassword = "需要修改密码"
)

type UserService interface {
	RegisterByEmailAndPassword(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	LoginByEmailAndPassword(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, string, error)
	GetUserInfoById(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	UpdateUserInfoById(ctx context.Context, userDomain domain.UserDomain) error
	DeleteUSerInfoByIds(ctx context.Context, ids []int64) error
	Logout(ctx context.Context, sessionId string) error
	LoginByPhoneCode(ctx context.Context, userDomain domain.UserDomain, biz string) (domain.UserDomain, error)
}

type UserServiceV1 struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceV1{
		userRepository: userRepository,
	}
}

func (u *UserServiceV1) LoginByPhoneCode(ctx context.Context, userDomain domain.UserDomain, biz string) (domain.UserDomain, error) {
	return u.userRepository.LoginOrRegister(ctx, userDomain)
}

func (u *UserServiceV1) Logout(ctx context.Context, sessionId string) error {
	return u.userRepository.Logout(ctx, sessionId)
}

func (u *UserServiceV1) RegisterByEmailAndPassword(ctx context.Context,
	userDomain domain.UserDomain) (domain.UserDomain, error) {
	// 密码加密
	gfp, err := bcrypt.GenerateFromPassword([]byte(userDomain.Password), bcrypt.DefaultCost)
	if err != nil {
		if err == ErrEmailExist {
			return domain.UserDomain{}, err
		}
		return domain.UserDomain{}, fmt.Errorf("加密失败，%w", err)
	}
	userDomain.Password = string(gfp)
	// 添加历史密码，新建用户不需要验证历史密码
	u.userRepository.SetHistoryPassword(ctx, userDomain.Email, userDomain.Password)
	// 创建用户
	return u.userRepository.CreateUser(ctx, userDomain)
}

func (u *UserServiceV1) LoginByEmailAndPassword(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, string, error) {
	user, err := u.userRepository.FindByEmail(ctx, userDomain)
	// 出现错误
	if err != nil {
		if err == ErrUserNotExist {
			return domain.UserDomain{}, "", err
		}
		return domain.UserDomain{}, "", err
	}
	// 密码比较
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDomain.Password))
	// 密码输入错误
	if err != nil {
		return domain.UserDomain{}, "", fmt.Errorf("密码输入错误，%w", err)
	}
	// 计算3个月的毫秒数（3个月 = 3 * 30天 * 24小时 * 60分钟 * 60秒 * 1000毫秒）
	threeMonthsMillis := int64(3 * 30 * 24 * 60 * 60 * 1000)
	// 获取当前时间的毫秒级时间戳
	now := time.Now().UnixMilli()
	// 判断是否需要修改密码
	if now-user.UpdateTime > threeMonthsMillis {
		return user, "需要修改密码", nil
	}
	return user, "", nil
}

func (u *UserServiceV1) GetUserInfoById(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	return u.userRepository.FindById(ctx, userDomain)
}

func (u *UserServiceV1) UpdateUserInfoById(ctx context.Context, userDomain domain.UserDomain) error {
	// 获取用户信息
	user, err := u.GetUserInfoById(ctx, userDomain)
	if err != nil {
		return err
	}
	// 密码解密，判断旧密码是否正确
	// 参数1：旧密码（加密后的），参数2：旧密码（明文）
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDomain.Password))
	if err != nil {
		// 旧密码输入错误
		return ErrPasswordWrong
	}
	// 获取历史密码
	historyPasswords, err := u.userRepository.GetHistoryPassword(ctx, user.Email)
	if err != nil {
		return err
	}
	// 判断新密码是否在历史密码库中
	for _, password := range historyPasswords {
		// 不相等，会返回err
		// 相等，会返回nil
		// 新密码在历史密码中
		if bcrypt.CompareHashAndPassword([]byte(password), []byte(userDomain.NewPassword)) == nil {
			return ErrUseOldPassword
		}
	}
	// 新密码加密
	gfp, er := bcrypt.GenerateFromPassword([]byte(userDomain.NewPassword), bcrypt.DefaultCost)
	if er != nil {
		return fmt.Errorf("加密失败，%w", err)
	}
	userDomain.Password = string(gfp)
	// 更新历史密码
	u.userRepository.SetHistoryPassword(ctx, user.Email, userDomain.Password)
	return u.userRepository.Update(ctx, userDomain)
}

func (u *UserServiceV1) DeleteUSerInfoByIds(ctx context.Context, ids []int64) error {
	return u.userRepository.DeleteByIds(ctx, ids)
}
