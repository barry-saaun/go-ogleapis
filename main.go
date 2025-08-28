package main

import (
	"fmt"
	"go-googleapis/auth"
	googleapispkg "go-googleapis/googleapis_pkg"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig *oauth2.Config

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variable")
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/tasks",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/calendar.events",
		},
		Endpoint: google.Endpoint,
	}

	// client, err := auth.getClient(oauthConfig)
	client, err := auth.GetClient(oauthConfig)
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
