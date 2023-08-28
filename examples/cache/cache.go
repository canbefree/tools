package main

import (
	"context"

	pb_example "github.com/SuperJourney/tools/examples/cache/proto"
	"github.com/SuperJourney/tools/helper"
	"golang.org/x/sync/singleflight"

	"github.com/golang/protobuf/proto"
)

type GrpcExampleService struct {
	Cache           *helper.Repo
	SingleTonFlight *singleflight.Group
}

func NewGrpcExampleService() *GrpcExampleService {
	return &GrpcExampleService{
		Cache:           helper.NewRepo(),
		SingleTonFlight: &singleflight.Group{},
	}
}

type Pb interface {
	pb_example.ExampleServer
}

const (
	Expired = 5
)

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

func (p *GrpcExampleService) hello(ctx context.Context, req *pb_example.HelloRequest) (*pb_example.HelloResponse, error) {
	return &pb_example.HelloResponse{}, nil
}
