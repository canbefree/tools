# tools

- 简单的函数类
- 基础库，比如说通用缓存，通用数据库，通用日志；


目录： 

helper : 简单的函数类；

infra: 基础库；

libs : 相关工具类封装；





- [X] gprc通用缓存
  ```
  func (p *GrpcExampleService) Hello(ctx context.Context, req *pb_example.HelloRequest) (resp *pb_example.HelloResponse, err error) {
  	var respWrap = &helper.GrpcRespWrap{
  		Message: resp,
  		Err:     nil,
  	}
  	if err := p.Cache.GrpcGet(ctx, p.SingleTonFlight, "Hello", req, respWrap, func() (proto.Message, error) {
  		return p.hello(ctx, req)
  	}, Expired); err != nil {
  		// 获取缓存失败直接查数据库
  		if err != respWrap.Err {
  			return p.hello(ctx, req)
  		}
  	}
  	return &pb_example.HelloResponse{}, nil
  }
  ```
- [ ] inner_event
