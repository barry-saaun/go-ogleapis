package googleapis_pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func ListTasks(client *http.Client) error {
	service, err := initTasksService(client)
	if err != nil {
		log.Fatalf("Unable to create Tasks Service: %w", err)
		return fmt.Errorf("Unable to create Tasks Service: %w", err)
	}

	taskLists, err := service.Tasklists.List().Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve task lists: %w", err)
	}

	fmt.Println("Your Task Lists:")
	if len(taskLists.Items) == 0 {
		fmt.Println("No task lists found. ✅")
	} else {
		for _, i := range taskLists.Items {
			fmt.Printf("- %s (%s)\n", i.Title, i.Id)
		}
	}

	return nil
}

func initTasksService(client *http.Client) (*tasks.Service, error) {
	ctx := context.Background()

	service, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return service, nil
}
