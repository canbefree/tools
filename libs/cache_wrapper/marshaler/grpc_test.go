package marshaler

import (
	"testing"

	demo "github.com/SuperJourney/tools/libs/cache_wrapper/marshaler/example"
	"github.com/stretchr/testify/assert"
)

func TestMarshaler_MarshalWrapper(t *testing.T) {
	type args struct {
		respBody []interface{}
	}
	tests := []struct {
		name      string
		m         *Marshaler
		args      args
		want      []byte
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test case 1",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					&demo.GetResponse{
						Data: "succ",
					},
					nil,
				},
			},
			want:      []byte(`{"Body":"CgRzdWNj","Err":null}`),
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.MarshalWrapper(tt.args.respBody...)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMarshaler_UnMarshalWrapper(t *testing.T) {
	var err error
	type args struct {
		strings []byte
		resp    any
	}
	tests := []struct {
		name      string
		m         *Marshaler
		args      args
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test case 1",
			m:    &Marshaler{},
			args: args{
				strings: []byte(`{"Body":"CgRzdWNj","Err":null}`),
				resp:    []interface{}{&demo.GetResponse{}, err},
			},
			want: []interface{}{&demo.GetResponse{
				Data: "succ",
			}, nil},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.UnMarshalWrapper(tt.args.strings, tt.args.resp)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
