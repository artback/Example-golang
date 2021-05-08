package status

import (
	"bitburst/pkg/online"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestClient_GetStatus(t *testing.T) {
	type args struct {
		ids []int
	}
	tests := []struct {
		name      string
		server    *httptest.Server
		closeConn bool
		args      args
		want      []online.Status
		wantErr   bool
	}{
		{
			name: "return status",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Send response to be tested
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args: args{[]int{1}},
			want: []online.Status{{Id: 1, Online: true}},
		},
		{
			name: "return empty",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte(`{}`))
			})),
			args:    args{ids: []int{1, 2}},
			wantErr: true,
		},
		{
			name: "return error status code",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args:    args{[]int{1, 4}},
			wantErr: true,
		},
		{
			name: "use baseURl",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			args:    args{[]int{1, 5}},
			wantErr: true,
		},
		{
			name: "fail get",
			server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte(`{ "online": true,"id": 1}`))
			})),
			closeConn: true,
			args:      args{[]int{1}},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closeConn {
				tt.server.Close()
			}
			c := NewClient(tt.server.Client(), tt.server.URL)
			got, err := c.GetStatus(tt.args.ids)
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
			want: NewClient(http.DefaultClient, "/"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.client, "/"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
