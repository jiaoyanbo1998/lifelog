package service

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"lifelog-grpc/feed/domain"
	"lifelog-grpc/feed/repository"
	"sort"
	"sync"
)

type FeedServiceV1 struct {
	feedRepository repository.FeedRepository
	// map[string]Handler key就是type，value就是具体的handler实现
	handlerMap map[string]Handler // 不同的业务有不同的handler实现
}

func NewFeedService(feedRepository repository.FeedRepository, handlerMap map[string]Handler) FeedService {
	return &FeedServiceV1{
		feedRepository: feedRepository,
		handlerMap:     handlerMap,
	}
}

// CreateFeedEvent 创建feed事件
func (f *FeedServiceV1) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 获取具体的handler
	handler, ok := f.handlerMap[feedEvent.Type]
	if !ok {
		return errors.New("代码错误，或业务方传递的feedEvent.Type错误")
	}
	// 调用具体的handler
	return handler.CreateFeedEvent(ctx, feedEvent)
}

func (f *FeedServiceV1) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	// errgroup.Group底层封装了sync.WaitGroup
	var eg errgroup.Group
	res := make([]domain.FeedEvent, 0, limit*int64(len(f.handlerMap)))
	// 读写锁
	var mu sync.RWMutex
	for _, handler := range f.handlerMap {
		// 使用局部变量，避免闭包问题
		h := handler
		// 并发调用
		eg.Go(func() error {
			feedEvents, err := h.FindFeedEvents(ctx, userId, createTime, limit)
			// 捕获错误
			if err != nil {
				return err
			}
			// 加锁
			mu.Lock()
			res = append(res, feedEvents...)
			// 解锁
			mu.Unlock()
			return nil
		})
	}
	// 等待所有goroutine执行完毕
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 聚合排序
	// 按照createTime从大到小排序
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreateTime > res[j].CreateTime
	})
	// 求limit和len(res)的最小值
	var n int
	if limit > int64(len(res)) {
		n = len(res)
	} else {
		n = int(limit)
	}
	// 返回前n个
	return res[:int(n)], nil
}
