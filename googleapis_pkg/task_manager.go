package googleapispkg

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
)

type TaskManager struct {
	TaskTitle string
	Note      *string
	Completed *string
	DueDate   string
}

type TaskManagerServices struct {
	tasksService    *tasks.Service
	calendarService *calendar.Service
}

type TaskResult struct {
	CreatedTask *tasks.Task
	taskId      string
	calendarId  string
}

func InitTaskAndCalendarService(client *http.Client) (*TaskManagerServices, error) {
	calendarSrv, calErr := initCalendarService(client)
	tasksSrv, taskErr := initTasksService(client)

	if calErr != nil || taskErr != nil {
		return nil, nil
	}

	return &TaskManagerServices{
		tasksService:    tasksSrv,
		calendarService: calendarSrv,
	}, nil
}

func CreateTask(taskManager *TaskManager, managerServices *TaskManagerServices, client *http.Client) (*TaskResult, error) {
	taskListId, err := GetTaskListId(managerServices.tasksService)
	if err != nil {
		return nil, err
	}

	task := &tasks.Task{
		Title:     taskManager.TaskTitle,
		Due:       taskManager.DueDate,
		Completed: nil,
	}

	createdTaskRes, err := managerServices.tasksService.Tasks.Insert(taskListId, task).Do()
	if err != nil {
		log.Fatalf("Failed creating task: %w\n", err)
	}

	fmt.Printf("Created Task ID: %s\n", createdTaskRes.Id)

	return &TaskResult{
		CreatedTask: createdTaskRes,
		taskId:      createdTaskRes.Id,
		calendarId:  "",
	}, nil
}

func ModifyTask() {}
