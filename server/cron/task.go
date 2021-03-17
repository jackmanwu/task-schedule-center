package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"task-schedule-center/db"
	"task-schedule-center/server/core"
	"task-schedule-center/service"
)

func InitTask(c *cron.Cron) {
	tasks, err := core.FindEnableTasks()
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		id, innerErr := c.AddJob(task.Cron, service.NewTaskJob(task.Id, task.Gid, task.Path, 0))
		if innerErr != nil {
			panic(err)
		}
		innerErr = db.RDB.Set(context.Background(), fmt.Sprintf(db.TaskInfo, task.Id), fmt.Sprintf("%d", id), 0).Err()
		if innerErr != nil {
			panic(innerErr)
		}
	}
}
