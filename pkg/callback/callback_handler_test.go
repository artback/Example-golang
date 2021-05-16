package callback_test

import (
	"bitburst/pkg/callback"
	"bitburst/pkg/id"
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type service struct {
	err error
}

func (s service) Handle(_ context.Context, _ []int) error {
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
			method: http.MethodPost,
			body:   []byte(`{"object_ids":[1,2]}`),
			want:   http.StatusOK,
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
		t.Run(tt.name, func(t *testing.T) {
			g := gin.Default()
			g.POST("/", callback.Handler(tt.service))
			req, _ := http.NewRequest(tt.method, "/", bytes.NewReader(tt.body))
			testHTTPResponse(g, req, func(w *httptest.ResponseRecorder) {
				if status := w.Code; status != tt.want {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tt.want)
				}
			})
		})
	}
}

// Helper function to process a request and test its response
func testHTTPResponse(r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder)) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	f(w)
}
