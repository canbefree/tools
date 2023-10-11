# tools

目录：

    helper : 简单的函数类；

    infra: 基础库；

    libs : 相关工具类封装；

## 相关工具类封装；

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
- [X] timewheel

  ```
  package main

  import (
	  "log"
	  "time"

	  timewheel "github.com/SuperJourney/tools/libs/time_wheel"
  )

  func main() {
	  timeWheel := timewheel.NewTimeWheel(time.Second, 10)
	  timeWheel.AddTask(10*time.Second, 1, func() {
		  log.Println("after 10 seconds, task 1 is done")
	  })
	  log.Println("time wheel start")
	  addTaskAfter10Seconds(timeWheel)
	  stopTimeWheelAfter30Seconds(timeWheel)
	  timeWheel.Start()

  }

  func addTaskAfter10Seconds(timeWheel *timewheel.TimeWheel) {
	  c := time.After(7 * time.Second)
	  go func() {
		  for range c {
			  timeWheel.AddTask(10*time.Second, 2, func() {
				  log.Println("after 17 seconds, task 2 is done")
			  })
		  }
	  }()
  }

  func stopTimeWheelAfter30Seconds(timeWheel *timewheel.TimeWheel) {
	  c := time.After(30 * time.Second)
	  go func() {
		  for range c {
			  log.Println("stop")
			  timeWheel.Stop()
			  return
		  }
	  }()
  }

  ```
