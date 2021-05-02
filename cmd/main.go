package main

import (
	"bitburst/internal/config"
	"bitburst/internal/logging"
	"bitburst/pkg/bitburst"
	"bitburst/pkg/online"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	repo, err := bitburst.NewPostgresRepository(conf.GetdbUrl())
	if err != nil {
		log.Fatal(err)
	}
	handler := bitburst.NewCallBackHandler(online.NewClient(client, conf.Service.Host), repo)

	go func() {
		for {
			func() {
				defer logging.Elapsed("PostgresRepository.DeleteOlder")()
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				err := repo.DeleteOlder(ctx, time.Now().Add(-30*time.Second))
				if err != nil {
					log.Println(fmt.Errorf("DeleteOlder %e", err))
				}
			}()
			time.Sleep(30 * time.Second)
		}
	}()
	router := mux.NewRouter()
	router.Handle("/callback", logging.Handler(handler)).Methods(http.MethodPost)
	srv := http.Server{
		Addr:         conf.Host,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		Handler:      router,
	}
	log.Println("Start server", conf.Host)
	log.Fatal(srv.ListenAndServe())
}
