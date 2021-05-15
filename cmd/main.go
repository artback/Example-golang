package main

import (
	"bitburst/internal/config"
	"bitburst/pkg/bitburst"
	"bitburst/pkg/bitburst/callback"
	"bitburst/pkg/bitburst/status"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	repo, err := bitburst.NewPostgresRepository(conf.GetdbUrl())
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			err := repo.DeleteOlder(ctx, time.Now().Add(-30*time.Second))
			cancel()
			if err != nil {
				log.Error(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	s := bitburst.Service{
		Client: status.NewClient(&http.Client{
			Timeout: time.Second * 5,
		}, conf.Service.Host),
		Repository: repo,
	}
	log.Info("Start server", conf.Host)
	log.Error(callback.SetupRouter(s).Run(conf.Host))
}
