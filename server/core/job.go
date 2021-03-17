package core

import (
	sql2 "database/sql"
	"fmt"
	"strings"
	"task-schedule-center/db"
)

func CreateJob(gid int, tid, taskDate int64, state string) (int64, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				fmt.Println(fmt.Sprintf("create job rollback err:%v", err))
			}
			panic(p)
		} else if err != nil {
			err = tx.Rollback()
			if err != nil {
				fmt.Println(fmt.Sprintf("create job rollback err:%v", err))
			}
		} else {
			err = tx.Commit()
			if err != nil {
				fmt.Println(fmt.Sprintf("create job commit tx err:%v", err))
			}
		}
	}()
	sql := `insert into job(gid,tid,task_date,create_time) values(?,?,?,unix_timestamp())`
	result, err := tx.Exec(sql, gid, tid, taskDate)
	if err != nil {
		return 0, err
	}
	jobId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	jobStateSql := `insert into job_state(job_id,state,time,create_time)values(?,?,unix_timestamp(),unix_timestamp())`
	_, err = tx.Exec(jobStateSql, jobId, state)
	if err != nil {
		return 0, err
	}
	return jobId, nil
}

func FindTaskJobCount(tids []string, taskDate int64, state string) (map[string]int, error) {
	sql := "select b.tid,count(1) as total from job_state a left join job b on a.job_id=b.id where b.tid in(" + strings.Join(tids, ",") + ") and a.state=? and b.task_date>=? and b.task_date <? group by b.tid"
	rows, err := db.DB.Query(sql, state, taskDate, taskDate+24*3600)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	result := map[string]int{}
	for rows.Next() {
		var tid string
		var count int
		err1 := rows.Scan(&tid, &count)
		if err1 != nil {
			return nil, err1
		}
		result[tid] = count
	}
	return result, nil
}

func SaveJobState(jobId int64, state string) (int64, error) {
	sql := "insert into job_state(job_id,state,time,create_time)values(?,?,unix_timestamp(),unix_timestamp())"
	result, err := db.DB.Exec(sql, jobId, state)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateJobWithIp(jobId, ip int64, state string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				fmt.Println(fmt.Sprintf("update job ip rollback err:%v", err))
			}
			panic(p)
		} else if err != nil {
			err = tx.Rollback()
			if err != nil {
				fmt.Println(fmt.Sprintf("update job ip rollback err:%v", err))
			}
		} else {
			err = tx.Commit()
			if err != nil {
				fmt.Println(fmt.Sprintf("create job commit tx err:%v", err))
			}
		}
	}()
	sql := `update job set ip=?,update_time=unix_timestamp() where id=?`
	result, err := tx.Exec(sql, ip, jobId)
	if err != nil {
		return err
	}
	affectRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectRows == 0 {
		return nil
	}
	jobStateSql := `insert into job_state(job_id,state,time,create_time)values(?,?,unix_timestamp(),unix_timestamp())`
	result, err = tx.Exec(jobStateSql, jobId, state)
	if err != nil {
		return err
	}
	affectRows, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if affectRows == 0 {
		err = sql2.ErrNoRows
		return err
	}
	return nil
}
