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

func (r *cacheWrapper) SetHandle(fn func(...interface{}) []interface{}) {
	r.handler = fn
}

func (r *cacheWrapper) Request(ctx context.Context,
	resp []interface{}, // UnMarshalWrapper 需要resp的类型
	reqs ...interface{},
) ([]interface{}, error) {
	if r.handler == nil {
		panic("handler not set yet")
	}
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
			curResp := r.handler(reqs...)
			respMarshal, err := r.MarshalWrapper(curResp)
			if err != nil {
				return nil, err
			}
			if err := r.cacheEngine.Set(ctx, key, respMarshal, int64(r.expired)); err != nil {
				return nil, err
			}
			return curResp, err
		}

		// 获取cache失败，直接访问原函数
		return r.handler(reqs...), nil
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

	handler func(...interface{}) []interface{}
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
