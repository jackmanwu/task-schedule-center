package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"task-schedule-center/db"
	"task-schedule-center/domain"
	"task-schedule-center/server/core"
	"task-schedule-center/util"
	"time"
)

type TaskJob struct {
	Tid      int64 `json:"tid"`
	JobId    int64 `json:"job_id"`
	gid      int
	Path     string `json:"path"`
	TaskDate int64  `json:"task_date"`
}

func NewTaskJob(tid int64, gid int, path string, taskDate int64) *TaskJob {
	return &TaskJob{Tid: tid, gid: gid, Path: path, TaskDate: taskDate}
}

func (t *TaskJob) Run() {
	AddJobToQueue(t)
}

func AddJobToQueue(t *TaskJob) int64 {
	if t.TaskDate == 0 {
		now, _ := util.ParseTime(time.Now().Format(util.YMD))
		t.TaskDate = now.Unix() - 24*3600
	}
	jobId, err := core.CreateJob(t.gid, t.Tid, t.TaskDate, domain.StateTrigger)
	if err != nil {
		fmt.Println(fmt.Sprintf("create job err: %v", err))
		return 0
	}
	//检验是否有任务已经在运行
	exists, err := existsJob(t.Tid, t.TaskDate)
	if err != nil {
		fmt.Println(fmt.Sprintf("task job exists err,jobId:%d,%v", jobId, err))
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}
	if exists {
		fmt.Println(fmt.Sprintf("task job exists,jobId:%d", jobId))
		saveJobState(jobId, domain.StateDuplicated)
		return jobId
	}
	//校验前置任务
	prevTid, err := getUncompletedPrevTid(t.Tid, t.TaskDate)
	if err != nil {
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}
	if prevTid > 0 {
		saveJobState(jobId, domain.StatePrevTaskUncompleted)

		waitLockKey := fmt.Sprintf(db.TaskWaitingSetLock, t.Tid, t.TaskDate)
		waitLock, err := db.RDB.SetNX(context.Background(), waitLockKey, prevTid, 10*time.Second).Result()
		if err != nil {
			return jobId
		}
		if waitLock {
			waitKey := fmt.Sprintf(db.TaskWaitingSet, prevTid, t.TaskDate)
			err = db.RDB.SAdd(context.Background(), waitKey, t.Tid).Err()
			if err != nil {
				fmt.Println(fmt.Sprintf("add wait tid err,tid:%d,prev_tid:%d,%v", t.Tid, prevTid, err))
			}
		}
		err = db.RDB.Del(context.Background(), waitLockKey).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("del wait set lock err,tid:%d,prev_tid:%d,%v", t.Tid, prevTid, err))
		}
		return jobId
	}
	ips, err := core.FindNodeByGid(t.gid)
	if err != nil {
		//未找到执行节点
		saveJobState(jobId, domain.StateNotFindNode)
		fmt.Println(fmt.Sprintf("not find node,tid:%d,%v", t.Tid, err))
		return jobId
	}

	//负载均衡，目前简单做成随机，后续根据任务数量、cpu、内存等节点负载情况做优化
	targetIp := ips[rand.Intn(len(ips))]
	key := fmt.Sprintf(db.TaskJobQueue, t.gid, targetIp)
	t.JobId = jobId
	value, err := json.Marshal(t)
	if err != nil {
		fmt.Println(fmt.Sprintf("serialization job err,job_id:%d,%v", jobId, err))
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}
	err = db.RDB.RPush(context.Background(), key, value).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("add job to queue err, job_id: %d,%v", jobId, err))
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}

	lockKey := fmt.Sprintf(db.TaskJobLock, t.Tid, t.TaskDate)
	err = db.RDB.Set(context.Background(), lockKey, jobId, 0).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("task job lock err, job_id: %d,%v", jobId, err))
		rollback(key, value, jobId)
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}
	err = core.UpdateJobWithIp(jobId, targetIp, domain.StateInQueue)
	if err != nil {
		fmt.Println(fmt.Sprintf("add job err,jobId:%d,%v", jobId, err))
		rollback(key, value, jobId)
		saveJobState(jobId, domain.StateInnerFailed)
		return jobId
	}

	fmt.Println(fmt.Sprintf("add job to queue success,job_id: %d", jobId))
	return jobId
}

func rollback(key string, value []byte, jobId int64) {
	//回滚
	err := db.RDB.LRem(context.Background(), key, 0, value).Err()
	if err != nil {
		fmt.Println(fmt.Sprintf("rollback redis queue err,jobId:%d,%v", jobId, err))
	}
}

func saveJobState(jobId int64, expectedState string) {
	affectRows, err := core.SaveJobState(jobId, expectedState)
	if err != nil || affectRows == 0 {
		fmt.Println(fmt.Sprintf("save job state err,job_id: %d,expected_state:%s,%v", jobId, expectedState, err))
	}
}

func getUncompletedPrevTid(tid int64, taskDate int64) (int64, error) {
	prevTids, err := core.FindPrevTask(tid)
	if err != nil {
		fmt.Println(fmt.Sprintf("find prev task err,tid:%d,%v", tid, err))
		return -1, err
	}
	if prevTids == nil {
		return 0, nil
	}
	var tids []string
	for _, prevTid := range prevTids {
		exists, err1 := existsJob(prevTid, taskDate)
		if err1 != nil {
			return -1, err1
		}
		if exists {
			return prevTid, nil
		}
		tids = append(tids, strconv.FormatInt(prevTid, 10))
	}
	if len(tids) == 0 {
		return 0, nil
	}
	tidJobCountMap, err := core.FindTaskJobCount(tids, taskDate, domain.StateCompleted)
	if err != nil {
		fmt.Println(fmt.Sprintf("When run job,find prev task count err,tid:%d, %v", tid, err))
		return -1, err
	}
	for _, ptid := range tids {
		if tidJobCountMap[ptid] == 0 {
			prevTid, _ := strconv.ParseInt(ptid, 10, 64)
			return prevTid, nil
		}
	}
	return 0, nil
}

func existsJob(tid int64, taskDate int64) (bool, error) {
	lockKey := fmt.Sprintf(db.TaskJobLock, tid, taskDate)
	val, err := db.RDB.Exists(context.Background(), lockKey).Result()
	if err != nil {
		return true, err
	}
	return val == 1, nil
}
