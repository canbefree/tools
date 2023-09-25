package rules

import (
	"context"
	"time"

	"github.com/SuperJourney/tools/helper"
)

type NatureWeek struct {
	BaseRule
}

func (n NatureWeek) Name() string {
	return "natureweek"
}

func (k *NatureWeek) Incr(ctx context.Context, cacheKey []byte) error {
	var exists bool
	var err error

	if exists, err = k.cache.Exist(ctx, cacheKey); err != nil {
		return err
	}
	// set new cache
	if !exists {
		var expireTime int64
		now := time.Now()
		lastTime := helper.GetLastWeekTime(now)
		expireTime = lastTime - now.Unix()

		if err := k.cache.Set(ctx, cacheKey, []byte("1"), expireTime); err != nil {
			return err
		}

		return nil
	}

	// update
	if err := k.cache.Incr(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}
