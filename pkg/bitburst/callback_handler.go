package bitburst

import (
	"bitburst/internal/logging"
	"bitburst/pkg/id"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type callbackHandler struct {
	id.Service
}

func NewCallBackHandler(service id.Service) http.Handler {
	return callbackHandler{service}
}

type response struct {
	ObjectIds []int `json:"object_ids"`
}

func (c callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		logging.Error.Println(fmt.Errorf("decode body %e", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		if err := c.Handle(ctx, resp.ObjectIds); err != nil {
			logging.Error.Println(err)
		}
	}()
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	return
}
