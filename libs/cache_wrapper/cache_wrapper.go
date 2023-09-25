package cache_wrapper

import (
	"context"

	"github.com/SuperJourney/tools/infra"
	"golang.org/x/sync/singleflight"
)

func NewCacheWrapper(cache infra.KVCacheI, request RequestFormatter, expired int) *cacheWrapper {
	return &cacheWrapper{
		cacheEngine:      cache,
		RequestFormatter: request,
		expired:          expired,
		group:            &singleflight.Group{},
	}
}

func (r *cacheWrapper) Get(ctx context.Context,
	directFn func(...interface{}) []interface{},
	resp []interface{}, // UnMarshalWrapper 需要resp的类型
	reqs ...interface{},
) ([]interface{}, error) {
	key := r.GetUniqKey(reqs...)
	var fn = func() (interface{}, error) {
		cache, err := r.cacheEngine.Get(ctx, key)
		if err == nil {
			resp, err := r.UnMarshalWrapper(cache, resp)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}

		if r.cacheEngine.IsNotFoundErr(err) {
			resp := directFn(reqs...)
			respMarshal, err := r.MarshalWrapper(resp)
			if err != nil {
				return nil, err
			}
			if err := r.cacheEngine.Set(ctx, key, respMarshal, int64(r.expired)); err != nil {
				return nil, err
			}
			return resp, err
		}

		// 获取cache失败，直接访问原函数
		return directFn(reqs...), nil
	}
	i, err, _ := r.group.Do(string(key), fn)
	if err != nil {
		return nil, err
	}

	ret := i.([]interface{})
	return ret, err
}

type cacheWrapper struct {
	cacheEngine infra.KVCacheI
	group       *singleflight.Group
	RequestFormatter
	expired int // 缓存多久
}

//go:generate mockery --name=RequestFormatter --outpkg=mock_formatter
type RequestFormatter interface {
	GetUniqKey(...interface{}) []byte
	MarshalWrapI
}

type MarshalWrapI interface {
	MarshalWrapper(...interface{}) ([]byte, error)
	UnMarshalWrapper([]byte, any) ([]interface{}, error)
}
