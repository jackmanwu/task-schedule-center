package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-schedule-center/web/common"
	"task-schedule-center/web/core"
)

type TaskBody struct {
	Gid  int    `form:"gid" binding:"required"`
	Name string `form:"name" binding:"required"`
	Cron string `form:"cron" binding:"required"`
	Path string `form:"path" binding:"required"`
	Uid  int64  `form:"uid" binding:"required"`
}

func CreateHandler(c *gin.Context) {
	var body TaskBody
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := core.Insert(body.Gid, body.Name, body.Cron, body.Path, body.Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessWithData(map[string]int64{"id": id}))
}
