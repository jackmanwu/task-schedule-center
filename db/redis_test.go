package db

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	//val, err := RDB.Exists(context.Background(), "snowflake:midset").Result()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(val)
	//
	//r, err := RDB.SRem(context.Background(), "wujinlei", "12").Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(r)

	//l, err := RDB.SetNX(context.Background(), "lock", 1, 0).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(l)
	var nexts []string
	nexts = append(nexts, "2")
	err := RDB.SRem(context.Background(), "tsc:w:s:1:1615824000", nexts).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func removeNextTids(tid int64, taskDate int64, waitKey string, nextTids *[]string) (bool, error) {
	waitLock, err := RDB.SetNX(context.Background(), TaskWaitingSetLock, tid, 10*time.Second).Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("get wait lock err,tid:%d,taskDate:%d,%v", tid, taskDate, err))
		return false, err
	}
	fmt.Println(waitLock)
	if waitLock {
		err = RDB.SRem(context.Background(), waitKey, nextTids).Err()
		if err != nil {
			fmt.Println(fmt.Sprintf("remove wait set err,tid:%d,%v", tid, err))
			return false, err
		}
		return true, nil
	}
	return false, nil
}
