package core

import "task-schedule-center/db"

func FindNodeByGid(gid int) ([]int64, error) {
	sql := `select ip from node where gid=? and state=1`
	rows, err := db.DB.Query(sql, gid)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var ips []int64
	for rows.Next() {
		var ip int64
		err = rows.Scan(&ip)
		ips = append(ips, ip)
	}
	return ips, nil
}
