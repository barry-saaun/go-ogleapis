package googleapispkg

import (
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
)

type TaskManager struct {
	TaskTitle string
	Note      *string
	Completed bool
	DueDate   time.Time
}

type TaskManagerServices struct {
	tasksService    *tasks.Service
	calendarService *calendar.Service
}

type TaskReult struct {
	CreatedTask *tasks.Task
	taskId      string
	talendarId  string
}

func initTaskAndCalendarService(client *http.Client) (*TaskManagerServices, error) {
	calendarSrv, calErr := InitCalendarService(client)
	tasksSrv, taskErr := InitTasksService(client)

	if calErr != nil || taskErr != nil {
		return nil, nil
	}

	return &TaskManagerServices{
		tasksService:    tasksSrv,
		calendarService: calendarSrv,
	}, nil
}

func CreateTask(taskManager *TaskManager, managerServices *TaskManagerServices, client *http.Client) (*TaskReult, error) {
	taskListId, err := GetTaskListId(managerServices.tasksService)
	if err != nil {
		return nil, err
	}
}

func ModifyTask() {}
