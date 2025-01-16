package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/lifeLog/domain"
	"lifelog-grpc/pkg/loggerx"
	"strconv"
	"time"
)

type LifeLogCache interface {
	GetFirstPage(ctx context.Context, authorId int64) ([]domain.LifeLogDomain, error)
	SetFirstPage(ctx context.Context, authorId int64, lifeLogDomains []domain.LifeLogDomain) error
	DelFirstPage(ctx context.Context, authorId int64) error
	Set(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error
	SetPublic(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error
	GetPublic(ctx context.Context, authorId int64) domain.LifeLogDomain
}

type LifeLogRedisCache struct {
	cmd    redis.Cmdable
	logger loggerx.Logger
}

func NewLifeLogRedisCache(cmd redis.Cmdable, l loggerx.Logger) LifeLogCache {
	return &LifeLogRedisCache{
		cmd:    cmd,
		logger: l,
	}
}

// GetFirstPage 获取redis中存储的LifeLog列表的第一页数据
func (a *LifeLogRedisCache) GetFirstPage(ctx context.Context, authorId int64) ([]domain.LifeLogDomain, error) {
	result, err := a.cmd.Get(ctx, a.GetFirstKey(authorId)).Result()
	if err != nil {
		a.logger.Error("redis获取失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:GetFirstPage"))
		return nil, err
	}
	// json反序列化，将[]byte类型的json数据转为domain.LifeLogDomain
	var lifeLogDomains []domain.LifeLogDomain
	err = json.Unmarshal([]byte(result), &lifeLogDomains)
	if err != nil {
		a.logger.Error("json反序列化失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:GetFirstPage"))
		return nil, err
	}
	return lifeLogDomains, nil
}

// Set 将LifeLog列表的第一页的第一条数据存储到redis中
func (a *LifeLogRedisCache) Set(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error {
	key := fmt.Sprintf("first_page_data:%d", lifeLogDomain.Id)
	// 序列化
	val, err := json.Marshal(lifeLogDomain)
	if err != nil {
		a.logger.Error("序列化失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Set"))
		return err
	}
	err = a.cmd.Set(ctx, key, val, time.Minute*1).Err()
	if err != nil {
		a.logger.Error("redis存储失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Set"))
		return err
	}
	return nil
}

// SetFirstPage 将LifeLog列表的第一页数据存储到redis中
func (a *LifeLogRedisCache) SetFirstPage(ctx context.Context, authorId int64, lifeLogDomains []domain.LifeLogDomain) error {
	// 将数据json序列化([]byte格式)
	//   为什么要将存储到redis中的数据进行json序列化？
	//		方便存储，可以跨语言使用，提高了查询和存储的效率
	bytes, err := json.Marshal(lifeLogDomains)
	if err != nil {
		a.logger.Error("序列化失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:SetFirstPage"))
		return err
	}
	// 将数据存储到redis
	key := a.GetFirstKey(authorId)
	err = a.cmd.Set(ctx, key, bytes, time.Minute*10).Err()
	if err != nil {
		a.logger.Error("redis存储失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:SetFirstPage"))
		return err
	}
	return nil
}

// DelFirstPage 删除redis中存储的LifeLog列表的第一页数据
func (a *LifeLogRedisCache) DelFirstPage(ctx context.Context, authorId int64) error {
	err := a.cmd.Del(ctx, a.GetFirstKey(authorId)).Err()
	if err != nil {
		a.logger.Error("删除缓存失败", loggerx.Error(err),
			loggerx.String("method:", "lifeLogRepository:Create"))
		return err
	}
	return nil
}

// SetLifeLogAbstract 设置LifeLog摘要
func (a *LifeLogRedisCache) SetLifeLogAbstract(ads []domain.LifeLogDomain) {
	for i := 0; i < len(ads); i++ {
		// 将内容替换为摘要
		runes := []rune(ads[i].Content)
		if len(runes) <= 100 {
			continue
		}
		ads[i].Content = string(runes[:100]) + "..."
	}
}

// GetFirstKey 获取key
func (a *LifeLogRedisCache) GetFirstKey(authorId int64) string {
	return "lifeLog_first_page_" + strconv.Itoa(int(authorId))
}

// SetPublic 将第一次发布的数据存储到redis中
func (a *LifeLogRedisCache) SetPublic(ctx context.Context, lifeLogDomain domain.LifeLogDomain) error {
	bytes, err := json.Marshal(lifeLogDomain)
	if err != nil {
		a.logger.Error("序列化失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:SetPublic"))
		return err
	}
	err = a.cmd.Set(ctx, a.GetPublicKey(lifeLogDomain.Author.Id), bytes, time.Minute*10).Err()
	if err != nil {
		a.logger.Error("redis存储失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:SetPublic"))
		return err
	}
	return nil
}

// GetPublic 从redis中获取第一次发布的数据
func (a *LifeLogRedisCache) GetPublic(ctx context.Context, authorId int64) domain.LifeLogDomain {
	result, err := a.cmd.Get(ctx, a.GetPublicKey(authorId)).Result()
	if err != nil {
		a.logger.Error("redis获取失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:GetPublic"))
		return domain.LifeLogDomain{}
	}
	var ad domain.LifeLogDomain
	err = json.Unmarshal([]byte(result), &ad)
	if err != nil {
		a.logger.Error("json反序列化失败", loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:GetPublic"))
		return domain.LifeLogDomain{}
	}
	return ad
}

// GetPublicKey 获取key
func (a *LifeLogRedisCache) GetPublicKey(authorId int64) string {
	return "lifeLog_public_" + strconv.Itoa(int(authorId))
}
