package bitburst

import (
	"bitburst/pkg/online"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type callbackHandler struct {
	online.Client
	online.Repository
	context context.Context
}

func NewCallBackHandler(ctx context.Context, client online.Client, repository online.Repository) http.Handler {
	return callbackHandler{Client: client, Repository: repository, context: ctx}
}

type response struct {
	ObjectIds []int `json:"object_ids"`
}

func (c callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
			err = c.UpsertAll(c.context, status)
			if err != nil {
				log.Println(err)
			}
		}()
		w.WriteHeader(http.StatusOK)
		return

	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
func readStatus(result  chan result) ([]online.Status, error) {
	var status []online.Status
	var err error
	for r := range result {
		if r.err != nil {
			err = r.err
		}
		if r.status != nil {
			status = append(status, *r.status)
		}
	}
	if err != nil {
		return nil, err
	}
	return status, nil
}
