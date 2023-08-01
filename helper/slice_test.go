package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type T any

func TestUnsetSlice(t *testing.T) {
	type args struct {
		s []T
		i int
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{
			name: "test del index 0 ",
			args: args{
				s: []T{1, 2, 3, 4, 5},
				i: 0,
			},
			want: []T{2, 3, 4, 5},
		},

		{
			name: "test del index out of range (large) ",
			args: args{
				s: []T{1, 2, 3, 4, 5},
				i: 5,
			},
			want: []T{1, 2, 3, 4, 5},
		},

		{
			name: "test del index out of range (less) ",
			args: args{
				s: []T{1, 2, 3, 4, 5},
				i: -1,
			},
			want: []T{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, UnsetSlice(tt.args.s, tt.args.i))
		})
	}
}
