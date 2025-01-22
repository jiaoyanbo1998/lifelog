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

//func NewFeedService(feedRepository repository.FeedRepository, handlerMap map[string]Handler) FeedService {
//	return &FeedServiceV1{
//		feedRepository: feedRepository,
//		handlerMap:     handlerMap,
//	}
//}

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

// FeedServiceV2 线程安全的FeedService
type FeedServiceV2 struct {
	feedRepository repository.FeedRepository
	handlers       sync.Map // 使用sync.Map替代原生map
}

func NewFeedService(feedRepository repository.FeedRepository, initialHandlers map[string]Handler) FeedService {
	svc := &FeedServiceV2{
		feedRepository: feedRepository,
	}
	// 初始化注册Handler
	for bizType, handler := range initialHandlers {
		svc.handlers.Store(bizType, handler)
	}
	return svc
}

// CreateFeedEvent 创建feed事件（线程安全版本）
func (f *FeedServiceV2) CreateFeedEvent(ctx context.Context, feedEvent domain.FeedEvent) error {
	// 使用Load替代map查找
	val, ok := f.handlers.Load(feedEvent.Type)
	if !ok {
		return errors.New("handler not found")
	}
	handler, ok := val.(Handler)
	if !ok {
		return errors.New("invalid handler type")
	}
	return handler.CreateFeedEvent(ctx, feedEvent)
}

// FindFeedEvents 查询优化版本
func (f *FeedServiceV2) FindFeedEvents(ctx context.Context, userId, createTime, limit int64) ([]domain.FeedEvent, error) {
	var (
		eg            errgroup.Group                      // 使用errgroup.Group，带有错误处理的sync.WaitGroup
		resCh         = make(chan []domain.FeedEvent, 10) // 缓冲通道提升性能
		merged        []domain.FeedEvent                  // 最终结果
		collectorDone = make(chan struct{})               // 信号量
	)
	// 结果收集协程
	go func() {
		// 关闭通道
		defer close(collectorDone)
		// 收集结果
		for feeds := range resCh {
			merged = append(merged, feeds...)
		}
	}()
	// 遍历所有Handler
	f.handlers.Range(func(key, value any) bool {
		// 类型断言
		h, ok := value.(Handler)
		if !ok {
			// 返回true表示继续遍历下一个元素
			return true // 跳过非法Handler
		}
		// 并发查询
		eg.Go(func() error {
			// 查询
			feedEvents, err := h.FindFeedEvents(ctx, userId, createTime, limit)
			if err != nil {
				return nil // 允许部分失败
			}
			// 阻塞当前协程
			select {
			// 发送结果
			case resCh <- feedEvents:
			// 等待ctx.Done()信号
			case <-ctx.Done():
				return ctx.Err() // 结束select
			}
			return nil // 允许部分失败
		})
		// 返回true表示继续遍历下一个元素
		return true
	})
	// 等待所有查询完成
	err := eg.Wait()
	if err != nil {
		close(resCh) // 关闭通道
		return nil, err
	}
	// 关闭通道
	close(resCh)
	<-collectorDone // 等待结果收集完成
	// 高效排序优化
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].CreateTime > merged[j].CreateTime
	})
	// 结果截取优化，只返回limit个结果
	if int64(len(merged)) > limit {
		return merged[:limit], nil
	}
	return merged, nil
}
