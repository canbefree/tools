package demo

import "context"

//go:generate mockery --name=DemoServiceI --outpkg=mock_demo

type DemoServiceI interface {
	Get(ctx context.Context, req *GetRequest) (*GetResponse, error)
}

// type DemoService struct {
// }

// func (s *DemoService) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
// 	return &GetResponse{Data: "succ"}, nil
// }
