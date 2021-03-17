package core

import (
	sql2 "database/sql"
	"task-schedule-center/db"
	"task-schedule-center/domain"
)

func FindEnableTasks() ([]domain.Task, error) {
	sql := `select id,gid,name,cron,path,create_time from task where state=1`
	rows, err := db.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Gid, &task.Name, &task.Cron, &task.Path, &task.CreateTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func FindPrevTask(tid int64) ([]int64, error) {
	sql := `select prev_tid from task_prev where tid=?`
	rows, err := db.DB.Query(sql, tid)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var tids []int64
	for rows.Next() {
		var prevTid int64
		if err = rows.Scan(&prevTid); err != nil {
			return nil, err
		}
		tids = append(tids, prevTid)
	}
	return tids, nil
}
