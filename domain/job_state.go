package domain

type JobState struct {
	Id         int64
	JobId      int64
	State      string
	Time       int64
	CreateTime int
}
