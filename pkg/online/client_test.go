package online_test

import (
	"bitburst/pkg/online"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetStatus(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name      string
		server    *httptest.Server
		closeConn bool
		args      args
		want      *online.Status
		wantErr   bool
	}{
		{
			name: "return status",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Send response to be tested
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args: args{id: 1},
			want: &online.Status{1, true},
		},
		{
			name: "return empty",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(`{}`))
			})),
			args:    args{id: 1},
			wantErr: true,
		},
		{
			name: "return error status code",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args:    args{id: 1},
			wantErr: true,
		},
		{
			name: "use baseURl",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args:    args{id: 1},
			wantErr: true,
		},
		{
			name: "fail get",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			closeConn: true,
			args:      args{id: 1},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closeConn {
				tt.server.Close()
			}
			c := online.NewClient(tt.server.Client(), tt.server.URL)
			got, err := c.GetStatus(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		client  *http.Client
		baseURL string
	}
	tests := []struct {
		name string
		args args
		want online.Client
	}{
		{
			name: "client",
			args: args{
				client:  http.DefaultClient,
				baseURL: "/",
			},
			want: online.NewClient(http.DefaultClient, "/"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := online.NewClient(tt.args.client, "/"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
