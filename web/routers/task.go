package routers

import (
	"github.com/gin-gonic/gin"
	"task-schedule-center/web/routers/handler"
)

func TaskRouters(e *gin.Engine) {
	v1 := e.Group("/v1/task")
	{
		v1.POST("/create", handler.CreateHandler)
	}
}
