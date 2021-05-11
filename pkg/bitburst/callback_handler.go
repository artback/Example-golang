package bitburst

import (
	"bitburst/pkg/id"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type response struct {
	ObjectIds []int `json:"object_ids"`
}

func CallBackHandler(service id.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp response
		err := c.ShouldBindJSON(&resp)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.Error(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			defer cancel()
			err := service.Handle(ctx, resp.ObjectIds)
			if err != nil {
				c.Error(err)
			}
		}()
	}
}
