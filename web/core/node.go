package core

import "task-schedule-center/db"

func UpdateNode(gid int, ip string) error {
	sql := `replace into node(gid,ip,state,update_time) values(?,inet_aton(?),1,unix_timestamp())`
	_, err := db.DB.Exec(sql, gid, ip)
	if err != nil {
		return err
	}
	return nil
}
