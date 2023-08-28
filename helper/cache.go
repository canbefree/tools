package helper

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/status"
	gstatus "google.golang.org/grpc/status"
	// "google.golang.org/protobuf/proto"
)

type GrpcRespWrap struct {
	proto.Message
	Err error `json:"err"`
}

type RespCache struct {
	Message    []byte         `json:"message,omitempty"`
	GrpcStatus *status.Status `json:"grpc_status,omitempty"`
	Err        error          `json:"err,omitempty"`
}

type CacheI interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key, value []byte, expired int) error
	IsNotFoundErr(error) bool
}

type Repo struct {
	cache CacheI
}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GrpcGet(ctx context.Context, group *singleflight.Group, scene string, msg proto.Message, resp *GrpcRespWrap, directFn func() (proto.Message, error), expire int) error {
	key, err := GetCommonProtoKey(scene, msg)
	if err != nil {
		return err
	}

	var getCache = func() (interface{}, error) {
		cache, err := r.cache.Get(ctx, key)
		if err == nil {
			var respCache *RespCache
			err := json.Unmarshal(cache, &respCache)
			if err != nil {
				return nil, fmt.Errorf("cache, json Unmarshal  err:%v", err)
			}
			if respCache.GrpcStatus != nil {
				resp.Err = gstatus.ErrorProto(respCache.GrpcStatus)
			} else {
				resp.Err = respCache.Err
			}

			if respCache.Message != nil {
				// 转换成proto信息
				err = proto.Unmarshal(respCache.Message, resp.Message)
				if err != nil {
					return nil, fmt.Errorf("cache, proto Unmarshal  err:%v", err)
				}

			}
			return resp.Message, resp.Err
		}
		// 缓存未读区，使用cache
		if r.cache.IsNotFoundErr(err) {
			respMsg, resperr := directFn()
			var respCacheMsg []byte
			if !proto.MessageV2(respMsg).ProtoReflect().IsValid() {
				respCacheMsg = nil
			} else {
				respCacheMsg, err = proto.Marshal(respMsg)
				if err != nil {
					return nil, fmt.Errorf("cache, proto Marshal err:%v", err)
				}
			}
			respCache := &RespCache{
				Message: respCacheMsg,
			}
			s, ok := gstatus.FromError(resperr)
			if !ok {
				respCache.Err = resperr
			} else {
				respCache.GrpcStatus = s.Proto()
			}

			respCacheByte, err := json.Marshal(respCache)
			if err != nil {
				return nil, fmt.Errorf("cache, json Marshal  err:%v", err)
			}
			// 设置缓存
			r.cache.Set(ctx, []byte(key), respCacheByte, expire)
			return respMsg, resperr
		}
		return nil, err
	}

	i, err, _ := group.Do(key, getCache)
	if err != nil {
		return err
	}
	if v, ok := i.(proto.Message); ok {
		resp.Message = v
		return nil
	}
	return fmt.Errorf("获取缓存失败")
}

func GetCommonProtoKey(scene string, r proto.Message) (string, error) {
	if scene == "" {
		return "", fmt.Errorf("scene empty")
	}
	if r == nil {
		return "", fmt.Errorf("message is nil")
	}
	pr, err := proto.Marshal(r)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write([]byte(pr))
	sig := hex.EncodeToString(m.Sum(nil))

	return fmt.Sprintf("common-cache:%v:%v", scene, sig), nil
}
