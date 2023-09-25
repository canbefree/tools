package cache_wrapper

import (
	"context"
	"errors"
	"testing"

	mock_infra_cache "github.com/SuperJourney/tools/infra/mocks"
	demo "github.com/SuperJourney/tools/libs/cache_wrapper/marshaler/example"
	mock_demo "github.com/SuperJourney/tools/libs/cache_wrapper/marshaler/example/mocks"
	mock_formatter "github.com/SuperJourney/tools/libs/cache_wrapper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_cacheWrapper_Get(t *testing.T) {
	var localCache []byte = nil

	cacheEngine := mock_infra_cache.NewKVCacheI(t)
	cacheEngine.On("Get", mock.Anything, mock.Anything).Return(func(ctx context.Context, key []byte) ([]byte, error) {
		if localCache != nil {
			return localCache, nil
		}
		return nil, errors.New("")
	})

	cacheEngine.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		if localCache == nil {
			localCache = args.Get(2).([]byte)
		}
	}).Return(nil)

	cacheEngine.On("IsNotFoundErr", mock.Anything).Return(func(err error) bool {
		return err != nil
	})

	grpcRequestFormatter := mock_formatter.NewRequestFormatter(t)
	grpcRequestFormatter.On("GetUniqKey", mock.Anything, mock.Anything).Return([]byte("uniqKey"))
	grpcRequestFormatter.On("MarshalWrapper", mock.Anything, mock.Anything).Return([]byte("resp-encode"), nil)
	grpcRequestFormatter.On("UnMarshalWrapper", mock.Anything, mock.Anything).Return([]interface{}{&demo.GetResponse{}, nil}, nil)

	request := mock_demo.NewDemoServiceI(t)
	request.On("Get", mock.Anything, mock.Anything).Return(&demo.GetResponse{}, nil)

	tests := []struct {
		name      string
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{}

	tests = append(tests, struct {
		name      string
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{
		name: "test-succ",
		want: []interface{}{&demo.GetResponse{}, nil},
		assertion: func(t assert.TestingT, err error, o ...interface{}) bool {
			return assert.NoError(t, err)
		},
	})

	tests = append(tests, struct {
		name      string
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{
		name: "test-cache",
		want: []interface{}{&demo.GetResponse{}, nil},
		assertion: func(t assert.TestingT, err error, o ...interface{}) bool {
			return assert.NoError(t, err)
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheWrapper := NewCacheWrapper(cacheEngine, grpcRequestFormatter, 10)

			var fn = func(args ...interface{}) []interface{} {
				resp, err := request.Get(args[0].(context.Context), args[1].(*demo.GetRequest))
				return []interface{}{resp, err}
			}
			got, err := cacheWrapper.Get(context.TODO(), fn, []interface{}{&demo.GetResponse{}, errors.New("")}, context.TODO(), &demo.GetRequest{})

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
