package status_test

import (
	"bitburst/pkg/online"
	"bitburst/pkg/status"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestDecodeStatus(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *online.Status
		wantErr bool
	}{
		{
			name: "decode Status body",
			args: args{
				r: strings.NewReader(`{ "online": true,"id": 1}`),
			},
			want: &online.Status{1, true},
		},
		{
			name: "decode nil",
			args: args{
				r: strings.NewReader(``),
			},
			wantErr: true,
		},
		{
			name: "decode no online",
			args: args{
				r: strings.NewReader(`{"id": 1}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := status.Decode(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
