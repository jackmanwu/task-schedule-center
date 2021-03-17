package core

import (
	sql2 "database/sql"
	"task-schedule-center/db"
	"task-schedule-center/domain"
)

func GetGroup(gid int) (*domain.Group, error) {
	sql := `select id,name,create_time from task_schedule_center.group where id=?`
	var group domain.Group
	err := db.DB.QueryRow(sql, gid).Scan(&group.Id, &group.Name, &group.CreateTime)
	if err != nil {
		if err == sql2.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}
