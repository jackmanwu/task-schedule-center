package routers

import (
	"github.com/gin-gonic/gin"
	"task-schedule-center/web/routers/handler"
)

func ApiRouters(e *gin.Engine) {
	v1 := e.Group("/v1/api")
	{
		v1.POST("/heartbeat", handler.HeartbeatHandler)
		v1.GET("/pull", handler.PullHandler)
		v1.POST("/report", handler.ReportHandler)
		v1.POST("/exec", handler.ExecHandler)
	}
}
