package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/pkg/bloomFilter"
	"lifelog-grpc/pkg/loggerx"
)

var (
	ErrCodeSendMany     = errors.New("短信验证码发送次数太多，达到上限明日再试")
	ErrCodeSendFrequent = errors.New("短信验证码发送太频繁")
	ErrVerifyToMany     = errors.New("验证码验证太频繁")
)

// 将lua脚本注入到UseLuaSetCode变量中
//go:embed lua/set_code.lua
var UseLuaSetCode string

// 将lua脚本注入到UseLuaVerifyCode变量中
//go:embed lua/verify_code.lua
var UseLuaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
	SetBlackPhone(ctx context.Context, phone string) error
	IsBlackPhone(ctx context.Context, phone string) (bool, error)
}

type CodeRedisCache struct {
	cmd         redis.Cmdable
	logger      loggerx.Logger
	pc          prometheus.Counter
	bloomFilter *bloomFilter.BloomFilter
}

func NewCodeCache(cmd redis.Cmdable, l loggerx.Logger, bloomFilter *bloomFilter.BloomFilter) CodeCache {
	// 创建一个计数器
	pc := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "jyb",                    // 命名空间
		Subsystem: "webook",                 // 子系统
		Name:      "sendToMany_code_count",  // 名称
		Help:      "发送短信达到上限的人数", // 帮助信息
	})
	// 将计数器注册到Prometheus
	prometheus.MustRegister(pc)
	return &CodeRedisCache{
		cmd:         cmd,
		logger:      l,
		pc:          pc,
		bloomFilter: bloomFilter,
	}
}

func (c *CodeRedisCache) IsBlackPhone(ctx context.Context, phone string) (bool, error) {
	res1, _ := c.bloomFilter.GetMurmur3BitMap(phone)
	res2, _ := c.bloomFilter.GetMD5BitMap(phone)
	res3, _ := c.bloomFilter.GetBLAKE2BitMap(phone)
	if res1 == 0 || res2 == 0 || res3 == 0 {
		// 不在黑名单
		return false, nil
	}
	// 在黑名单
	return true, nil
}

func (c *CodeRedisCache) SetBlackPhone(ctx context.Context, phone string) error {
	err := c.bloomFilter.SetMurmur3BitMap(phone)
	if err != nil {
		return err
	}
	err = c.bloomFilter.SetMD5BitMap(phone)
	if err != nil {
		return err
	}
	err = c.bloomFilter.SetBLAKE2BitMap(phone)
	if err != nil {
		return err
	}
	return nil
}

func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, code string) error {
	// 验证输入参数
	if biz == "" || phone == "" || code == "" {
		return fmt.Errorf("参数不能为空")
	}
	key := c.GetKey(biz, phone)
	// 使用lua脚本，验证验证码
	res, err := c.cmd.Eval(ctx, UseLuaVerifyCode, []string{key}, code).Int()
	if err != nil {
		c.logger.Error("验证码验证失败", loggerx.Error(err))
		return fmt.Errorf("验证码验证失败，%w", err)
	}
	switch res {
	case -1:
		return fmt.Errorf("用户输入错误，%w", err)
	case 0:
		return nil
	case -2:
		return ErrVerifyToMany
	default:
		return fmt.Errorf("未知错误，%w", err)
	}
}

func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	// 验证输入参数
	if biz == "" || phone == "" || code == "" {
		return fmt.Errorf("参数不能为空")
	}
	// 获取key
	key := c.GetKey(biz, phone)
	// 使用lua脚本，存储验证码
	res, err := c.cmd.Eval(ctx, UseLuaSetCode, []string{key}, code).Int()
	if err != nil {
		c.logger.Error("验证码存储失败", loggerx.Error(err))
		return err
	}
	switch res {
	// 系统错误，不允许发短信
	case -1:
		return errors.New("系统错误")
	// 验证码发送成功
	case 0:
		return nil
	// 验证码发送太频繁
	case -2:
		return ErrCodeSendFrequent
	// 短信发送次数达到上限
	case -3:
		c.logger.Error("短信"+
			"发送次数达到上限", loggerx.String("phone", phone))
		c.pc.Inc() // 增加计数器
		return ErrCodeSendMany
	// 未知错误
	default:
		return errors.New("未知错误")
	}
}

func (c *CodeRedisCache) GetKey(biz, phone string) string {
	return fmt.Sprintf("%s:%s:code", biz, phone)
}
