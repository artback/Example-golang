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
	easy "github.com/t-tomalak/logrus-easy-formatter"
	ginlogrus "github.com/toorop/gin-logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}
	ginLogger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %method%(%statusCode%): %path% \n",
		},
	}
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
	router.Use(ginlogrus.Logger(ginLogger))
	router.POST("/callback", callback.Handler(s))
	logger.Error(router.Run(conf.ServerAddress))
}
