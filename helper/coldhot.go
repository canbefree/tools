package helper

import "context"

type CodeHotIFace interface {
	Get(ctx context.Context) ([]byte, error)
	Set(ctx context.Context, value []byte) error
}

type CodeHot struct {
	hot  map[string][]byte
	cold map[string][]byte
}

var _ CodeHotIFace = &CodeHot{}

func (c *CodeHot) Get(ctx context.Context) ([]byte, error) {
	// TODO 冷热交换实现，高并发场景下使用；
	panic("implement me")
}

func (c *CodeHot) Set(ctx context.Context, value []byte) error {
	// TODO 冷热交换实现，高并发场景下使用

	return nil
}
