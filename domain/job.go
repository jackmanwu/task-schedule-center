package domain

const (
	StateTrigger             = "TRIGGER"
	StateDuplicated          = "DUPLICATED"
	StatePrevTaskUncompleted = "PREV_TASK_UNCOMPLETED"
	StateNotFindNode         = "NOT_FOUND_NODE"
	StateInQueue             = "IN_QUEUE"
	StateGotIt               = "GOT_IT"
	StateStart               = "START"
	StateCompleted           = "COMPLETED"
	StateFailed              = "FAILED"
	StateInnerFailed         = "INNER_FAILED"
)

type Job struct {
	Id         int64
	Gid        int
	Tid        int64
	Ip         int64
	TaskDate   int64
	Param      string
	CreateTime int
	UpdateTime int
}
