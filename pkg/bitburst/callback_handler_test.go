package bitburst

import (
	"bitburst/pkg/online"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type client struct {
	status online.Status
	error  error
}

func (c client) GetStatus(_ int) (*online.Status, error) {
	return &c.status, c.error
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

func Test_callbackHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name       string
		client     online.Client
		repository online.Repository
		method     string
		body       []byte
		want       int
	}{
		{
			name: "ok request",
			client: client{
				status: *online.NewStatus(1, true),
			},
			repository: testRepository{},
			method:     http.MethodPost,
			body:       []byte(`{"object_ids":[1,2]}`),
			want:       http.StatusOK,
		},
		{
			name: "error response",
			client: client{
				error:  errors.New("something"),
				status: *online.NewStatus(1, true),
			},
			repository: testRepository{},
			body:       []byte(`{"object_ids":[1,2]}`),
			method:     http.MethodPost,
			want:       http.StatusOK,
		},
		{
			name: "error request",
			client: client{
				error:  errors.New("something"),
				status: *online.NewStatus(1, true),
			},
			repository: testRepository{},
			body:       []byte(`{"object_ids":["1","2"]}`),
			method:     http.MethodPost,
			want:       http.StatusBadRequest,
		},
		{
			name:   "wrong method",
			body:   []byte(`{"object_ids":[1,2]}`),
			method: http.MethodPatch,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name: "errors repository",
			client: client{
				status: *online.NewStatus(1, true),
			},
			repository: testRepository{
				err: errors.New("something"),
			},
			method: http.MethodPost,
			body:   []byte(`{"object_ids":[1,2]}`),
			want:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, "/", bytes.NewReader(tt.body))
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := NewCallBackHandler(context.Background(), tt.client, tt.repository)
			c.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.want)
			}
		})
	}
}

func TestNewCallBackHandler(t *testing.T) {
	tests := []struct {
		name string
		online.Client
		online.Repository
		want http.Handler
	}{
		{
			name:       "create callBackHandler",
			Client:     online.NewClient(http.DefaultClient, "/"),
			Repository: testRepository{},
			want:       NewCallBackHandler(context.Background(), online.NewClient(http.DefaultClient, "/"), testRepository{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCallBackHandler(context.Background(), tt.Client, tt.Repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCallBackHandler() = %v, want %v", got, tt.want)
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
					{status: online.NewStatus(1, true)},
					{status: online.NewStatus(2, true)},
				},
			},
			want: []online.Status{*online.NewStatus(1, true), *online.NewStatus(2, true)},
		},
		{
			name: "return error",
			args: args{
				input: []result{},
			},
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
