package core

import (
	"fmt"
	"strings"
	"task-schedule-center/db"
	"task-schedule-center/domain"
)

func SaveJobState(jobStates []*domain.JobState) (int64, error) {
	base := "insert into job_state(job_id,state,time,create_time)values"
	sub := "(%d,%s,%d,unix_timestamp())"
	var subs []string
	for _, jobState := range jobStates {
		subs = append(subs, fmt.Sprintf(sub, jobState.JobId, "'"+jobState.State+"'", jobState.Time))
	}
	all := base + strings.Join(subs, ",")
	result, err := db.DB.Exec(all)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
