package bitburst

import (
	"bitburst/pkg/online"
	"errors"
	"reflect"
	"testing"
)

func Test_getResult(t *testing.T) {
	tests := []struct {
		name    string
		client  online.Client
		ids     []int
		want    []online.Status
		wantErr bool
	}{
		{
			name: "return 1 status",
			client: client{
				status: online.Status{1, true},
			},
			ids:  []int{1},
			want: []online.Status{{1, true}},
		},
		{
			name: "return error",
			client: client{
				status: online.Status{1, true},
				error:  errors.New("something"),
			},
			ids:     []int{1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := getResult(tt.ids, tt.client)
			var got []online.Status
			var err error
			for r := range c {
				if r.err != nil {
					err = r.err
					break
				}
				got = append(got, *r.status)

			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readStatus(t *testing.T) {
	type args struct {
		input []result
	}
	tests := []struct {
		name    string
		args    args
		want    []online.Status
		wantErr bool
	}{
		{
			name: "all succeed",
			args: args{
				input: []result{
					{status: &online.Status{1, true}},
					{status: &online.Status{2, true}},
				},
			},
			want: []online.Status{{1, true}, {2, true}},
		},
		{
			name: "return error",
			args: args{
				input: []result{
					{status: nil, err: errors.New("something happened")},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := make(chan result)
			go func() {
				for _, input := range tt.args.input {
					r <- input
				}
				close(r)
			}()
			got, err := readStatus(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
