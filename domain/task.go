package domain

const (
	StateDisable = 0
	StateEnable  = 1
)

type Task struct {
	Id         int64
	Gid        int
	Name       string
	Cron       string
	State      uint8
	Path       string
	Uid        int64
	CreateTime int
	UpdateTime int
}
