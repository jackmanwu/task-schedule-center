package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"task-schedule-center/db"
	"task-schedule-center/domain"
	"task-schedule-center/service"
	"task-schedule-center/util"
	"task-schedule-center/web/common"
	"task-schedule-center/web/core"
	"time"
)

type Param struct {
	Gid int    `form:"gid" binding:"required"`
	Ip  string `form:"ip" binding:"required"`
}

type TaskJob struct {
	Tid      int64  `json:"tid"`
	JobId    int64  `json:"job_id"`
	Path     string `json:"path"`
	TaskDate int64  `json:"task_date"`
}

type Report struct {
	Tid      int64  `json:"tid"`
	JobId    int64  `json:"job_id"`
	State    string `json:"state"`
	Time     int64  `json:"time"`
	TaskDate int64  `json:"task_date"`
}

type ReportParams struct {
	Report []*Report `json:"report"`
}

func HeartbeatHandler(c *gin.Context) {
	var param Param
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if !validateGid(param.Gid, c) {
		return
	}
	if err := core.UpdateNode(param.Gid, param.Ip); err != nil {
		fmt.Println(fmt.Sprintf("When heartbeat err,gid: %d,ip:%s, %v", param.Gid, param.Ip, err))
	}
	c.JSON(http.StatusOK, common.NewSuccess())
}

func PullHandler(c *gin.Context) {
	var param Param
	if err := c.ShouldBind(&param); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if !validateGid(param.Gid, c) {
		return
	}
	num, err := util.Ip2Num(param.Ip)
	if err != nil {
		fmt.Println(fmt.Sprintf("When pull,convert ip to num err,gid:%d,ip:%s\n%v", param.Gid, param.Ip, err))
		return
	}
	key := fmt.Sprintf(db.TaskJobQueue, param.Gid, num)
	result, err := db.RDB.LPop(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusOK, common.NewSuccess())
			return
		}
		fmt.Println(fmt.Sprintf("When pull,pop err: %v", err))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	var taskJob TaskJob
	err = json.Unmarshal([]byte(result), &taskJob)
	if err != nil {
		fmt.Println(fmt.Sprintf("When pull,deserialization err: %v", err))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessWithData(taskJob))
}

func ReportHandler(c *gin.Context) {
	var params ReportParams
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	if len(params.Report) == 0 {
		return
	}
	var jobStates []*domain.JobState
	var triggerNext bool
	for _, report := range params.Report {
		jobStates = append(jobStates, &domain.JobState{JobId: report.JobId, State: report.State, Time: report.Time})
		if report.State == domain.StateCompleted {
			triggerNext = true
		}
	}
	affectRows, err := core.SaveJobState(jobStates)
	if err != nil {
		fmt.Println(fmt.Sprintf("When save job state,err: %v", err))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if affectRows == 0 {
		fmt.Println(fmt.Sprintf("When save job affect 0 rows"))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if triggerNext {
		tid := params.Report[0].Tid
		taskDate := params.Report[0].TaskDate
		go afterClean(tid, taskDate)
		go triggerNextTask(tid, taskDate)
	}
	c.JSON(http.StatusOK, common.NewSuccess())
}

func ExecHandler(c *gin.Context) {
	tidStr := c.PostForm("tid")
	tid, err := strconv.ParseInt(tidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "tid must be number")
		return
	}
	date := c.PostForm("task_date")
	tm, err := util.ParseTime(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, "task_date must format yyyy-MM-dd")
		return
	}
	taskDate := tm.Unix()
	task, err := core.GetTask(tid)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, common.NewResult(1, "task not exists", nil))
			return
		}
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if task.State == domain.StateDisable {
		c.JSON(http.StatusOK, common.NewResult(2, "task has been disabled", nil))
		return
	}
	jobId := service.AddJobToQueue(service.NewTaskJob(tid, task.Gid, task.Path, taskDate))
	if jobId == 0 {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if jobId == -1 {
		c.JSON(http.StatusOK, common.NewResult(3, "task is running", nil))
		return
	}
	c.JSON(http.StatusOK, common.NewSuccessWithData(map[string]int64{"job_id": jobId}))
}

func validateGid(gid int, c *gin.Context) bool {
	group, err := core.GetGroup(gid)
	if err != nil {
		fmt.Println(fmt.Sprintf("When find group err,gid:%d,%v", gid, err))
		c.JSON(http.StatusInternalServerError, nil)
		return false
	}
	if group == nil {
		c.JSON(http.StatusOK, common.NewResult(1, "group not exists", nil))
		return false
	}
	return true
}

func triggerNextTask(tid int64, taskDate int64) {
	waitKey := fmt.Sprintf(db.TaskWaitingSet, tid, taskDate)
	nextTids, err := db.RDB.SMembers(context.Background(), waitKey).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("get next tids err,tid:%d,%v", tid, err))
		return
	}
	if len(nextTids) == 0 {
		return
	}

	tasks, err := core.GetTaskBatch(nextTids)
	if err != nil {
		fmt.Println(fmt.Sprintf("trigger next task err:%v", err))
		return
	}
	if tasks == nil {
		return
	}
	for _, task := range *tasks {
		jobTask := service.NewTaskJob(task.Id, task.Gid, task.Path, taskDate)
		go service.AddJobToQueue(jobTask)
	}
	for i := 0; i < 10; i++ {
		r, err := removeNextTids(tid, taskDate, waitKey, &nextTids)
		if err != nil {
			continue
		}
		if r {
			break
		}
	}
}

func removeNextTids(tid int64, taskDate int64, waitKey string, nextTids *[]string) (bool, error) {
	waitLockKey := fmt.Sprintf(db.TaskWaitingSetLock, tid, taskDate)
	waitLock, err := db.RDB.SetNX(context.Background(), waitLockKey, tid, 10*time.Second).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("get wait lock err,tid:%d,taskDate:%d,%v", tid, taskDate, err))
		return false, err
	}
	defer func() {
		err = db.RDB.Del(context.Background(), waitLockKey).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("del wait lock err,tid:%d,task_date:%d,%v", tid, taskDate, err))
		}
	}()
	if waitLock {
		err = db.RDB.SRem(context.Background(), waitKey, *nextTids).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("remove wait set err,tid:%d,%v", tid, err))
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func afterClean(tid, taskDate int64) {
	lockKey := fmt.Sprintf(db.TaskJobLock, tid, taskDate)
	for i := 0; i < 5; i++ {
		err := db.RDB.Del(context.Background(), lockKey).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("del lock key err,tid:%d,taskDate:%d,%v", tid, taskDate, err))
			continue
		}
		break
	}
}
