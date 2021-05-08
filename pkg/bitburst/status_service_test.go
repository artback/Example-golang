package bitburst

import (
	"bitburst/pkg/online"
	"context"
	"errors"
	"testing"
	"time"
)

type client struct {
	status []online.Status
	error  error
}

func (c client) GetStatus(_ []int) ([]online.Status, error) {
	return c.status, c.error
}

type testRepository struct {
	err error
}

func (t testRepository) UpsertAll(_ context.Context, _ []online.Status, _ time.Time) error {
	return t.err
}
func (t testRepository) DeleteOlder(_ context.Context, _ time.Time) error {
	return t.err
}

func TestService_Handle(t *testing.T) {
	type fields struct {
		Client     online.Client
		Repository online.Repository
	}
	type args struct {
		ctx context.Context
		ids []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "errors client",
			fields: fields{
				Client: client{
					status: []online.Status{{Id: 1, Online: true}},
					error:  errors.New("shit happens"),
				},
				Repository: testRepository{},
			},
			args:    args{ctx: context.Background(), ids: []int{1, 2}},
			wantErr: true,
		},
		{
			name: "errors repository",
			fields: fields{
				Client: client{
					status: []online.Status{{Id: 1, Online: true}},
				},
				Repository: testRepository{
					err: errors.New("something"),
				},
			},
			args:    args{ctx: context.Background(), ids: []int{1, 2}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Client:     tt.fields.Client,
				Repository: tt.fields.Repository,
			}
			if err := s.Handle(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
