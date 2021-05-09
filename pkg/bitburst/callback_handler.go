package bitburst

import (
	"bitburst/internal/body"
	"bitburst/pkg/id"
	"context"
	log "github.com/sirupsen/logrus"
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
	status, err := body.DecodeJSONBody(r, &resp)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), status)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go func() {
		defer cancel()
		if err := c.Handle(ctx, resp.ObjectIds); err != nil {
			log.Error(err)
		}
	}()
	w.Header().Set("Connection", "close")
	return
}
