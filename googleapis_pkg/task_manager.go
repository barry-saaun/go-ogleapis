package googleapispkg

import "time"

type TaskManager struct {
	TaskId     string
	CalendarId string
	TaskTitle  string
	Note       *string
	Completed  bool
	DueDate    time.Time
}
