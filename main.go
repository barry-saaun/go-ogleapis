package main

import (
	"fmt"
	"go-googleapis/auth"
	googleapispkg "go-googleapis/googleapis_pkg"
	"log"

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

	taskManager := &googleapispkg.TaskManager{
		TaskTitle: "hello world new",
		DueDate:   "2025-08-29T23:00:00.000Z",
	}

	services, err := googleapispkg.InitTaskAndCalendarService(client)
	if err != nil {
		log.Fatalf("Error initialising services for task and calendar: %w\n", err)
	}

	createdTask, err := googleapispkg.CreateTask(taskManager, services, client)

	fmt.Printf("createdTaskId: %w\n", createdTask)

	fmt.Println("Done ✅")
}
