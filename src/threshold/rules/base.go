package rules

import (
	"context"

	"github.com/SuperJourney/tools/infra"
)

type I interface {
	Name() string
	Incr(ctx context.Context, cacheKey string) error
}

type BaseRule struct {
	cache infra.KVCacheI
}
