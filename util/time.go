package util

import "time"

const (
	YMD = "2006-01-02"
)

func FormatYMD(date int64) string {
	tm := time.Unix(date, 0)
	return tm.Format(YMD)
}

func ParseTime(date string) (time.Time, error) {
	return time.ParseInLocation(YMD, date, time.Local)
}
