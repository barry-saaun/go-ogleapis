package googleapispkg

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func ListTasks(service *tasks.Service, taskListId string) ([]*tasks.Task, error) {
	tasksRes, err := service.Tasks.List(taskListId).Do()
	if err != nil {
		return nil, fmt.Errorf("❌ Unable to retrieve your tasks: %w", err)
	}

	if len(tasksRes.Items) == 0 {
		return []*tasks.Task{}, nil
	}

	return tasksRes.Items, nil
}

func GetTaskListId(service *tasks.Service) (string, error) {
	taskLists, err := service.Tasklists.List().Do()
	if err != nil {
		return "", fmt.Errorf("❌ Unable to retrieve the task list: %w\n", err)
	}

	taskListId := taskLists.Items[0].Id

	return taskListId, nil
}

func InitTasksService(client *http.Client) (*tasks.Service, error) {
	ctx := context.Background()

	service, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return service, nil
}
