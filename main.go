package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
	tokenChan   = make(chan *oauth2.Token)
)

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
			"https://www.googleapis.com/auth/tasks.readonly",
		},
		Endpoint: google.Endpoint,
	}

	http.HandleFunc("/oauth2callback", handleCallback)

	go func() {
		port := "6769"
		fmt.Println("Server started at http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":6769", nil))
	}()

	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Login URL: ", url)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Open this URL in your browser? (y/n): ")
	answer, _ := reader.ReadString('\n')

	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		err := browser.OpenURL(url)
		if err != nil {
			log.Printf("Failed to open browser automatically: %v", err)
			fmt.Println("Please open the URL manually:", url)
		}
	} else {
		fmt.Println("Please open the URL manually")
	}

	token := <-tokenChan
	log.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Println("✅ Authentication successful!")
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed in exchange code", http.StatusInternalServerError)
	}

	tokenChan <- token

	fmt.Fprint(w, "✅ Authentication Successful! You can close this window.")
}
