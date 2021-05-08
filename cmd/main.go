package main

import (
	"bitburst/internal/config"
	"bitburst/internal/logging"
	"bitburst/pkg/bitburst"
	"bitburst/pkg/status"
	"context"
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
	repo, err := bitburst.NewPostgresRepository(conf.GetdbUrl())
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			err := repo.DeleteOlder(ctx, time.Now().Add(-30*time.Second))
			cancel()
			if err != nil {
				logging.Error.Println(err)
			}
			time.Sleep(30 * time.Second)

		}
	}()

	if err != nil {
		logging.Error.Fatal(err)
	}
	s := bitburst.Service{
		Client: status.NewClient(&http.Client{
			Timeout: time.Second * 5,
		}, conf.Service.Host),
		Repository: repo,
	}
	router := mux.NewRouter()
	router.Handle("/callback", logging.Handler(bitburst.NewCallBackHandler(s))).Methods(http.MethodPost)
	srv := http.Server{
		Addr:         conf.Host,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
		Handler:      router,
	}
	logging.Info.Println("Start server", conf.Host)
	logging.Error.Fatal(srv.ListenAndServe())
}
