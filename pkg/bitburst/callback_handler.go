package bitburst

import (
	"bitburst/pkg/online"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type callbackHandler struct {
	online.Client
	online.Repository
}

func NewCallBackHandler(client online.Client, repository online.Repository) http.Handler {
	return callbackHandler{Client: client, Repository: repository}
}

type response struct {
	ObjectIds []int `json:"object_ids"`
}

func (c callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		result := getResult(resp.ObjectIds, c.Client)
		status, err := readStatus(result)
		if err != nil {
			log.Println(err)
		}
		err = c.UpsertAll(context.Background(), status, time.Now())
		if err != nil {
			log.Println(err)
		}
	}()
	w.WriteHeader(http.StatusOK)
	return
}
