package main

import (
	"context"
	"fmt"
	"go-googleapis/auth"
	googleapispkg "go-googleapis/googleapis_pkg"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variable")
	}

	client, err := auth.GetClient()
	if err != nil {
		log.Fatalf("Unable to get client: %v\n", err)
	}

	ctx := context.Background()

	taskManager := &googleapispkg.TaskManager{
		TaskTitle: "test",
		// DueTime:   time.Date(2025, time.August, 29, 12, 0, 0, 0, time.Local),
		DueTime: time.Date(2025, time.August, 31, 14, 30, 0, 0, time.Local),
	}

	services, err := googleapispkg.InitTaskAndCalendarService(client)
	if err != nil {
		log.Fatalf("Error initialising services for task and calendar: %w\n", err)
	}

	createdTask, err := googleapispkg.CreateTaskWithDueTime(ctx, taskManager, services, client)
	if err != nil {
		log.Fatalf("❌ Failed to create task: %v", err)
	}
	fmt.Printf("createdTaskId: %s\n", createdTask.TaskId)

	fmt.Println("Done ✅")
}
