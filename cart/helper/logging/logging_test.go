package logging

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	type args struct {
		isDebug bool
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Logger
	}{
		{
			name: "Test logging success",
			args: args{
				isDebug: false,
			},
			want: nil,
		},
		{
			name: "Test logging failed",
			args: args{
				isDebug: false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.isDebug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
