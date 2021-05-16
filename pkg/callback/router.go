package callback

import (
	"bitburst/pkg/id"
	"github.com/gin-gonic/gin"
)

func SetupRouter(s id.Service) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger())
	e.POST("/callback", Handler(s))
	return e
}
