package main

import (
	"bitburst/internal/config"
	"bitburst/internal/logging"
	"bitburst/pkg/bitburst"
	"bitburst/pkg/online"
	"context"
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
	handler := bitburst.NewCallBackHandler(context.Background(), online.NewClient(client, conf.Service.Host), repo)

	go func() {
		for {
			repo.DeleteOlder(context.Background(), time.Now().Add(-30*time.Second))
			time.Sleep(30 * time.Second)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/callback", logging.Handler(handler))
	srv := http.Server{
		Addr:         conf.Host,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		Handler:      mux,
	}
	log.Println("Start server", conf.Host)
	log.Fatal(srv.ListenAndServe())
}
