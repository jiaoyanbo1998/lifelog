package ioc

import (
	"github.com/redis/go-redis/v9"
	"lifelog-grpc/pkg/bloomFilter"
)

func InitBloomFilter(cmd redis.Cmdable) *bloomFilter.BloomFilter {
	filter := bloomFilter.NewBloomFilter(
		cmd,
		"user:black:phone",
	)
	return filter
}
