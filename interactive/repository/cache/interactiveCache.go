package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/pkg/loggerx"
)

// 注入lua脚本
//go:embed lua/interactive.lua
var luaInteractiveScript string

type InteractiveCache interface {
	InsertReadCount(ctx context.Context, biz string, bizId int64) error
	InsertLikeCount(ctx context.Context, biz string, bizId int64) error
	DecreaseLikeCount(ctx context.Context, biz string, bizId int64) error
	InsertCollectCount(ctx context.Context, biz string, bizId int64) error
	DecreaseCollectCount(ctx context.Context, biz string, bizId int64) error
}

type InteractiveRedisCache struct {
	cmd    redis.Cmdable
	logger loggerx.Logger
}

func NewInteractiveCache(cmd redis.Cmdable, l loggerx.Logger) InteractiveCache {
	return &InteractiveRedisCache{
		cmd:    cmd,
		logger: l,
	}
}

func (i *InteractiveRedisCache) InsertReadCount(ctx context.Context, biz string, bizId int64) error {
	key := i.GetKey(biz, bizId)
	// 插入点赞数，执行lua脚本
	res, err := i.cmd.Eval(ctx, luaInteractiveScript, []string{key}, "read_count", 1).Int()
	if err != nil {
		i.logger.Error("lua脚本执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveCache:InsertReadCount"))
		return err
	}
	// 插入失败
	if res == 0 {
		i.logger.Warn("阅读数自增失败，缓存过期，或者有人在搞你的系统",
			loggerx.String("method:", "InteractiveCache:InsertReadCount"))
		return errors.New("阅读数自增失败，key不存在，缓存过期，或者有人在搞你的系统")
	}
	// if res == 1 插入成功
	return nil
}

func (i *InteractiveRedisCache) InsertLikeCount(ctx context.Context, biz string, bizId int64) error {
	key := i.GetKey(biz, bizId)
	// 插入点赞数，执行lua脚本
	res, err := i.cmd.Eval(ctx, luaInteractiveScript, []string{key}, "Like_count", 1).Int()
	if err != nil {
		i.logger.Error("lua脚本执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveCache:InsertLikeCount"))
		return err
	}
	// 插入失败
	if res == 0 {
		i.logger.Warn("点赞数自增失败，缓存过期，或者有人在搞你的系统",
			loggerx.String("method:", "InteractiveCache:InsertLikeCount"))
		return errors.New("点赞数自增失败，key不存在，缓存过期，或者有人在搞你的系统")
	}
	// if res == 1 插入成功
	return nil
}

// DecreaseLikeCount 减少点赞数
func (i *InteractiveRedisCache) DecreaseLikeCount(ctx context.Context, biz string, bizId int64) error {
	// 获取key
	key := i.GetKey(biz, bizId)
	res, err := i.cmd.Eval(ctx, luaInteractiveScript, []string{key}, -1).Int()
	if err != nil {
		i.logger.Error("lua脚本执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveCache:DecreaseLikeCount"))
		return err
	}
	if res == 0 {
		i.logger.Warn("点赞数自减失败，缓存过期，或者有人在搞你的系统", loggerx.String("method:", "InteractiveCache:DecreaseLikeCount"))
		return errors.New("点赞数自减失败，缓存过期，或者有人在搞你的系统")
	}
	return nil
}

// InsertCollectCount 插入收藏数
func (i *InteractiveRedisCache) InsertCollectCount(ctx context.Context, biz string, bizId int64) error {
	// 获取key
	key := i.GetKey(biz, bizId)
	// 插入redis，使用lua脚本
	res, err := i.cmd.Eval(ctx, luaInteractiveScript, []string{key}, "collect_count", 1).Int()
	// 插入失败
	if err != nil {
		i.logger.Error("lua脚本执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveCache:InsertCollectCount"))
		return err
	}
	// 收藏数自增失败
	if res == 0 {
		i.logger.Error("收藏数自增失败，key不存在，缓存过期，或者有人在搞你的系统",
			loggerx.String("method:", "InteractiveCache:InsertCollectCount"))
		return errors.New("收藏数自增失败，key不存在，缓存过期，或者有人在搞你的系统")
	}
	return nil
}

// DecreaseCollectCount 减少收藏数
func (i *InteractiveRedisCache) DecreaseCollectCount(ctx context.Context, biz string, bizId int64) error {
	// 获取key
	key := i.GetKey(biz, bizId)
	// 执行lua脚本
	res, err := i.cmd.Eval(ctx, luaInteractiveScript, []string{key}, "collect_count", -1).Int()
	if err != nil {
		i.logger.Error("lua脚本执行失败", loggerx.Error(err),
			loggerx.String("method:", "InteractiveCache:DecreaseCollectCount"))
		return err
	}
	if res == 0 {
		i.logger.Error("收藏数自减失败，key不存在，缓存过期，或者有人在攻击你的系统",
			loggerx.String("method:", "InteractiveCache:DecreaseCollectCount"))
		return errors.New("收藏数自减失败，key不存在，缓存过期，或者有人在攻击你的系统")
	}
	return nil
}

func (i *InteractiveRedisCache) GetKey(biz string, bizId int64) string {
	return fmt.Sprintf("%s:%d:interactive", biz, bizId)
}
