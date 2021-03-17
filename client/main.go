package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os/exec"
	"task-schedule-center/domain"
	"task-schedule-center/util"
	"time"
)

const (
	Host = "http://127.0.0.1:8080"
)

type TaskJob struct {
	Tid      int64  `json:"tid"`
	JobId    int64  `json:"job_id"`
	Path     string `json:"path"`
	TaskDate int64  `json:"task_date"`
}

type Result struct {
	Code int8    `json:"code"`
	Msg  string  `json:"msg"`
	Data TaskJob `json:"data"`
}

type Report struct {
	Tid      int64  `json:"tid"`
	JobId    int64  `json:"job_id"`
	State    string `json:"state"`
	Time     int64  `json:"time"`
	TaskDate int64  `json:"task_date"`
}

func NewReport(tid int64, jobId int64, state string, time int64, taskDate int64) *Report {
	return &Report{Tid: tid, JobId: jobId, State: state, Time: time, TaskDate: taskDate}
}

func main() {
	go func() {
		for {
			heartbeat()
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		go pull()
		time.Sleep(10 * time.Second)
	}
}

func heartbeat() {
	url := Host + "/v1/api/heartbeat"
	params := url2.Values{}
	params.Add("gid", "1")
	params.Add("ip", "127.0.0.1")
	resp, err := http.PostForm(url, params)
	if err != nil {
		fmt.Println(fmt.Sprintf("When heartbeat post err,gid:%d,ip:%s: %v", 1, "127.0.0.1", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("When heartbeat post failed,gid:%d,ip:%s,status_code:%d", 1, "127.0.0.1", resp.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("When heartbeat read err,gid:%d,ip:%s: %v", 1, "127.0.0.1", err))
		return
	}
	fmt.Println(fmt.Sprintf("heartbeat result: %s", string(body)))
}

func pull() {
	url := Host + "/v1/api/pull?gid=1&ip=127.0.0.1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(fmt.Sprintf("When pull,get err:%v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("When pull,get failed,status_code:%d", resp.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("When pull,read err:%v", err))
		return
	}
	fmt.Println(fmt.Sprintf("pull result: %s", string(body)))

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(fmt.Sprintf("When pull,deserialization err:%v", err))
		return
	}
	handleMsg(&result)
}

func handleMsg(result *Result) {
	var reports []*Report
	defer func(report []*Report) {
		if len(reports) > 0 {
			go reportJobState(reports)
		}
	}(reports)
	if result.Code == 0 {
		if result.Data.JobId == 0 {
			fmt.Println("When pull,get job null")
			return
		}
		tid := result.Data.Tid
		jobId := result.Data.JobId
		taskDate := result.Data.TaskDate
		reports = append(reports, NewReport(tid, jobId, domain.StateGotIt, time.Now().Unix(), taskDate))
		reports = append(reports, NewReport(tid, jobId, domain.StateStart, time.Now().Unix(), taskDate))
		cmd := exec.Command("/bin/bash", result.Data.Path, util.FormatYMD(result.Data.TaskDate))
		fmt.Println(cmd.Args)
		output, err := cmd.Output()
		if err != nil {
			reports = append(reports, NewReport(tid, jobId, domain.StateFailed, time.Now().Unix(), taskDate))
			fmt.Println(fmt.Sprintf("When exec shell err,shell:%s,%v", result.Data.Path, err))
			return
		}
		fmt.Println(fmt.Sprintf("job exec finished,shell:%s,result:%s", result.Data.Path, string(output)))
		reports = append(reports, NewReport(tid, jobId, domain.StateCompleted, time.Now().Unix(), taskDate))
	} else {
		fmt.Println(fmt.Sprintf("pull invalid code,code:%d", result.Code))
	}
}

func reportJobState(reports []*Report) {
	url := Host + "/v1/api/report"
	value, err := json.Marshal(map[string]interface{}{"report": reports})
	if err != nil {
		fmt.Println(fmt.Sprintf("report serialization err:%v", err))
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(value))
	if err != nil {
		fmt.Println(fmt.Sprintf("report post err,value:%s,%v", string(value), err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(fmt.Sprintf("report post failed,value:%s,status_code:%d", string(value), resp.StatusCode))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("report read err,value:%s,err:%v", string(value), err))
		return
	}
	fmt.Println(fmt.Sprintf("report result: %s", string(body)))
}
