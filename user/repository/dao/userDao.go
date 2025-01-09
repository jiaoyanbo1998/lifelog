package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"lifelog-grpc/errs"
	"lifelog-grpc/user/domain"
	"strings"
	"time"
)

type UserDao interface {
	Insert(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error)
	GetUserById(ctx context.Context, id int64) (domain.UserDomain, error)
	UpdateById(ctx context.Context, user User) error
	DeleteByIds(ctx context.Context, ids []int64) error
	GetUserByPhoneOrInsert(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error)
	UpdateAvatar(ctx context.Context, userId int64, filePath string) error
}

type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{
		db: db,
	}
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// sql.NullString
	// 		type NullString struct {
	//			String string
	//			Valid  bool
	//			当String == null时，Valid == false，当String != null时，Valid == true
	//		}
	Email         sql.NullString `gorm:"uniqueIndex"`
	Password      string
	CreateTime    int64
	UpdateTime    int64
	Phone         sql.NullString `gorm:"uniqueIndex"`
	WechatUnionId string
	WechatOpenId  sql.NullString `gorm:"uniqueIndex"`
	NickName      string
	Authority     int64
	Avatar        string
}

func (User) TableName() string {
	return "tb_user"
}

func (g *GormUserDao) UpdateAvatar(ctx context.Context, userId int64, filePath string) error {
	tx := g.db.WithContext(ctx).Model(&User{}).Where("id = ?", userId).
		Updates(map[string]any{
			"avatar":      filePath,
			"update_time": time.Now().UnixMilli(),
		})
	if tx.RowsAffected == 0 {
		return errors.New("数据库更新失败")
	}
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// GetUserByPhoneOrInsert 通过手机号查询用户信息，如果没有查到就插入
func (g *GormUserDao) GetUserByPhoneOrInsert(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	var user User
	return userDomain, g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询到数据err == nil，
		// 没有查询到数据，err == "record not found"，
		// 查询语句执行失败，返回err == 错误信息
		err := g.db.Where("phone = ?", userDomain.Phone).First(&user).Error
		// 查询操作出错
		if err != nil {
			// 没有查询到数据
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ud, er := g.Insert(ctx, userDomain)
				// 插入失败
				if er != nil {
					userDomain = domain.UserDomain{}
					return errors.New("数据库插入失败")
				}
				// 插入成功
				userDomain = ud
				return nil
			}
			// 查询语句执行失败
			userDomain = domain.UserDomain{}
			// 查询语句执行错误
			return errors.New("数据库查询失败")
		}
		// 用户已经注册成功
		userDomain = domain.UserDomain{
			Id:        user.Id,
			Email:     user.Email.String,
			Phone:     user.Phone.String,
			NickName:  user.NickName,
			Authority: user.Authority,
		}
		return nil
	})
}

// Insert 创建用户
func (g *GormUserDao) Insert(ctx context.Context, userDomain domain.UserDomain) (domain.UserDomain, error) {
	var user User
	now := time.Now().UnixMilli()
	user.Email = sql.NullString{
		String: userDomain.Email,
		Valid:  userDomain.Email != "", // 有值时条件成立为true，无值时条件不成立为false
	}
	user.Password = userDomain.Password
	user.CreateTime = now
	user.UpdateTime = now
	user.Phone = sql.NullString{
		String: userDomain.Phone,
		Valid:  userDomain.Phone != "",
	}
	user.WechatUnionId = userDomain.WechatUnionId
	user.WechatOpenId = sql.NullString{
		String: userDomain.WechatOpenId,
		Valid:  userDomain.WechatOpenId != "",
	}
	// 给用户名设个默认值
	name := uuid.New().String()
	user.NickName = name
	user.Authority = 2
	err := g.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		ok := strings.Contains(err.Error(), "Duplicate")
		if ok {
			return domain.UserDomain{}, errs.EmailExist
		}
		return domain.UserDomain{}, fmt.Errorf("数据库插入失败，%w", err)
	}
	return domain.UserDomain{}, nil
}

// GetUserByEmail 根据邮箱查询用户
func (g *GormUserDao) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	var user User
	// 不管有没有查到，都会返回nil，除非数据库查询错误
	err := g.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		// 用户没有被注册，数据库中没有数据
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserDomain{}, errs.UserNotExist
		}
		return domain.UserDomain{}, fmt.Errorf("数据库查询失败，%w", err)
	}
	return domain.UserDomain{
		Id:         user.Id,
		Email:      user.Email.String,
		Password:   user.Password,
		NickName:   user.NickName,
		UpdateTime: user.UpdateTime,
		Authority:  user.Authority,
	}, nil
}

// GetUserById 根据id查询用户
func (g *GormUserDao) GetUserById(ctx context.Context, id int64) (domain.UserDomain, error) {
	var user User
	// 不管有没有查到，都会返回nil，除非数据库查询错误
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		// 用户没有被注册，数据库中没有数据
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.UserDomain{}, errs.UserNotExist
		}
		return domain.UserDomain{}, fmt.Errorf("数据库查询失败，%w", err)
	}
	return domain.UserDomain{
		Id:       user.Id,
		Email:    user.Email.String,
		Phone:    user.Phone.String,
		NickName: user.NickName,
		Password: user.Password,
		Avatar:   user.Avatar,
	}, nil
}

func (g *GormUserDao) UpdateById(ctx context.Context, user User) error {
	err := g.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", user.Id).Updates(user).Error // 忽略零值(0、""、false)，修改其他字段
	if err != nil {
		return fmt.Errorf("数据库更新失败，%w", err)
	}
	return nil
}

func (g *GormUserDao) DeleteByIds(ctx context.Context, ids []int64) error {
	res := g.db.WithContext(ctx).Where("id in ?", ids).Delete(&User{})
	if res.Error != nil {
		return fmt.Errorf("数据库删除失败，%w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("要删除的数据不存在，%w", res.Error)
	}
	return nil
}
