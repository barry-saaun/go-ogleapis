package main

import (
	"fmt"
	"go-googleapis/auth"
	"go-googleapis/googleapis_pkg"
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

	// err = googleapis.ListTasks(client)
	err = googleapis_pkg.ListTasks(client)
	if err != nil {
		log.Fatalf("Error listting tasks: %v\n", err)
	}

	fmt.Println("Done ✅")
}
