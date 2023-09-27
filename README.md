# tools

目录：

    helper : 简单的函数类；

    infra: 基础库；

    libs : 相关工具类封装；



- [X] 通用缓存用例；

  ```golang

  import (
    	marshaler "github.com/SuperJourney/grpc_marshaler"
  )

  type RequestFormat struct {
  	marshaler.Marshaler
  }

  func NewRequestFormat() *RequestFormat {
  	return &RequestFormat{}
  }

  func (r *RequestFormat) GetUniqKey(key string, reqs ...interface{}) []byte {
  	bts, _ := r.MarshalWrapper(reqs...)
  	md5bts := md5.Sum(bts)
  	return []byte(fmt.Sprintf("common-cache:%s:%s", key, md5bts))
  }
  ```

  ```golang
  	cacheWrapper := NewCacheWrapper(cacheEngine, RequestFormat(), 10)
  	cacheWrapper.SetHandle(func(args ...interface{}) []interface{} {
  		resp, err := request.Get(args[0].(context.Context), args[1].(*DemoRequest))
  		return []interface{}{resp, err}
  	})
  	got, err := cacheWrapper.Request(context.TODO(), []interface{}{&DemoResponse{}, errors.New("")}, context.TODO(), tt.demoRequest)

  ```


- [X] inner_event

    MQ平替版
