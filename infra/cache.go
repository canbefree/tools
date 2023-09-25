package infra

import "context"

type BaseCache interface {
	Exist(ctx context.Context, key []byte) (bool, error)
}

//go:generate mockery --name=KVCacheI --outpkg=mock_infra_cache
type KVCacheI interface {
	BaseCache
	Get(ctx context.Context, key []byte) ([]byte, error)
	Set(ctx context.Context, key, value []byte, expired int64) error
	Incr(ctx context.Context, key []byte) error
	IsNotFoundErr(error) bool
}

type DefaultCache struct {
}

func (c *DefaultCache) Get(ctx context.Context, key string) ([]byte, error) {
	panic("implement me")
}
func (c *DefaultCache) Set(ctx context.Context, key, value []byte, expired int) error {
	panic("implement me")
}
func (c *DefaultCache) IsNotFoundErr(err error) bool {
	panic("implement me")
}
