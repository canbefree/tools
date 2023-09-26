package cache_wrapper

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_infra_cache "github.com/SuperJourney/tools/infra/mocks"
	mock_formatter "github.com/SuperJourney/tools/libs/cache_wrapper/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DemoService struct {
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

type DemoRequest struct {
	Id int64
}

type DemoResponse struct {
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

func (s *DemoService) Get(ctx context.Context, req *DemoRequest) (*DemoResponse, error) {
	if req.Id == 0 {
		return nil, fmt.Errorf("id is 0")
	}
	return &DemoResponse{
		Name: "name",
		Age:  1,
	}, nil
}

var localCache []byte = nil

func Test_cacheWrapper_Get(t *testing.T) {

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
	grpcRequestFormatter.On("MarshalWrapper", mock.Anything, mock.Anything).Return([]byte(`{"name": "name", "age": 1}`), nil)
	grpcRequestFormatter.On("UnMarshalWrapper", mock.Anything, mock.Anything).Return(func([]byte, interface{}) ([]interface{}, error) {
		return []interface{}{&DemoResponse{
			Name: "name",
			Age:  1,
		}, nil}, nil
	}).Maybe()

	request := NewDemoService()

	tests := []struct {
		name        string
		demoRequest *DemoRequest
		want        []interface{}
		assertion   assert.ErrorAssertionFunc
	}{}

	tests = append(tests, struct {
		name        string
		demoRequest *DemoRequest
		want        []interface{}
		assertion   assert.ErrorAssertionFunc
	}{
		name: "test-succ",
		demoRequest: &DemoRequest{
			Id: 1,
		},
		want: []interface{}{&DemoResponse{
			Name: "name",
			Age:  1,
		}, nil},
		assertion: func(t assert.TestingT, err error, o ...interface{}) bool {
			return assert.NoError(t, err)
		},
	})

	tests = append(tests, struct {
		name        string
		demoRequest *DemoRequest
		want        []interface{}
		assertion   assert.ErrorAssertionFunc
	}{
		name: "test-cache",
		demoRequest: &DemoRequest{
			Id: 1,
		},
		want: []interface{}{&DemoResponse{
			Name: "name",
			Age:  1,
		}, nil},
		assertion: func(t assert.TestingT, err error, o ...interface{}) bool {
			return assert.NoError(t, err)
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheWrapper := NewCacheWrapper(cacheEngine, grpcRequestFormatter, 10)
			var fn = func(args ...interface{}) []interface{} {
				resp, err := request.Get(args[0].(context.Context), args[1].(*DemoRequest))
				return []interface{}{resp, err}
			}
			got, err := cacheWrapper.Get(context.TODO(), fn, []interface{}{&DemoResponse{}, errors.New("")}, context.TODO(), tt.demoRequest)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
