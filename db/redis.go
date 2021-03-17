package db

import "github.com/go-redis/redis/v8"

const (
	TaskInfo           = "tsc:t:i:%d"
	TaskJobQueue       = "tsc:j:q:%d:%d"
	TaskJobLock        = "tsc:t:j:l:%d:%d"
	TaskWaitingSet     = "tsc:w:s:%d:%d"
	TaskWaitingSetLock = "tsc:w:s:l:%d:%d"
)

var RDB = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})
