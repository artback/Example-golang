package bitburst

import (
	"bitburst/pkg/id"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type service struct {
	err error
}

func (s service) Handle(_ context.Context, ids []int) error {
	return s.err
}
func Test_callbackHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name    string
		service id.Service
		method  string
		body    []byte
		want    int
	}{
		{
			name: "ok request",
			service: service{
				nil,
			},
			body: []byte(`{"object_ids":[1,2]}`),
			want: http.StatusOK,
		},
		{
			name: "error response",
			service: service{
				nil,
			},
			body:   []byte(`{"object_ids":[1,2]}`),
			method: http.MethodPost,
			want:   http.StatusOK,
		},
		{
			name: "error request",
			service: service{
				nil,
			},
			body:   []byte(`{"object_ids":["1","2"]}`),
			method: http.MethodPost,
			want:   http.StatusBadRequest,
		},
		{
			name: "error service",
			service: service{
				errors.New("shit happens"),
			},
			body:   []byte(`{"object_ids":[1,2]}`),
			method: http.MethodPost,
			want:   http.StatusOK,
		},
	}
	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, "/", bytes.NewReader(tt.body))
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := NewCallBackHandler(tt.service)
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
		name    string
		service id.Service
		want    http.Handler
	}{
		{
			name: "create callBackHandler",
			want: NewCallBackHandler(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCallBackHandler(tt.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCallBackHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
