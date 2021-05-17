package main

import (
	"bitburst/internal/pkg/config"
	"bitburst/internal/pkg/elapsed"
	"bitburst/pkg/bitburst"
	"bitburst/pkg/callback"
	"bitburst/pkg/status"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"net/http"
	"time"
)

func main() {
	logger := logrus.New()
	conf, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	repo, err := bitburst.NewPostgresRepository(conf.DBSource)
	if err != nil {
		logger.Fatal(err)
		return
	}

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			err := repo.DeleteOlder(ctx, time.Now().Add(-30*time.Second))
			cancel()
			if err != nil {
				logger.Error(err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	s := bitburst.Service{
		Client: status.NewClient(&http.Client{}, conf.ServiceAddress),
		Upsert: repo,
		Log:    elapsed.Log{Info: logger},
	}

	logger.Info("Start server ", conf.ServerAddress)
	router := gin.New()
	router.Use(ginlogrus.Logger(logger))
	router.POST("/callback", callback.Handler(s))
	logger.Error(router.Run(conf.ServerAddress))
}
