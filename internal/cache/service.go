package cache

import (
	"context"
	"errors"
	"log"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key string) (val []byte, found bool, err error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration, keepTTL bool) (evicted bool, err error)
	Delete(ctx context.Context, key string) (deleted bool, err error)
	Len(ctx context.Context) (uint64, error)
}

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrTooLarge        = errors.New("too large")
)

type CacheServiceImpl struct {

}

func NewCacheServiceImpl() CacheService {
	return &CacheServiceImpl{

	}
}

func (s *CacheServiceImpl) Get(ctx context.Context, key string) (val []byte, found bool, err error) {
	result := []byte("Get requested")
	log.Printf("Got get request: key: %v\n", key)

	return result, true, nil
}

func (s *CacheServiceImpl) Set(ctx context.Context, key string, val []byte, ttl time.Duration, keepTTL bool) (evicted bool, err error) {
	var result bool
	log.Printf("Got set request: key: %v, val: %v, ttl: %v\n", key, val, ttl)

	return result, nil
}


func (s *CacheServiceImpl) Delete(ctx context.Context, key string) (deleted bool, err error) {
	return false, nil
}


func (s *CacheServiceImpl) Len(ctx context.Context) (uint64, error) {
	return 100, nil
}