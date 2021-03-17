package core

import (
	sql2 "database/sql"
	"strings"
	"task-schedule-center/db"
	"task-schedule-center/domain"
)

func Insert(gid int, name, cron string, path string, uid int64) (int64, error) {
	sql := `insert into task(gid,name,cron,path,uid,create_time) values(?,?,?,?,?,unix_timestamp())`
	result, err := db.DB.Exec(sql, gid, name, cron, path, uid)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func GetTaskBatch(tids []string) (*[]domain.Task, error) {
	sql := "select id,gid,path from task where id in (" + strings.Join(tids, ",") + ") and state=1"
	rows, err := db.DB.Query(sql)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Gid, &task.Path)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return &tasks, err
}

func GetTask(tid int64) (*domain.Task, error) {
	sql := `select id,gid,path,state from task where id=?`
	var task domain.Task
	err := db.DB.QueryRow(sql, tid).Scan(&task.Id, &task.Gid, &task.Path, &task.State)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
