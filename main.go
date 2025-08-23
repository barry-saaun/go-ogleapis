package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
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
			"https://www.googleapis.com/auth/tasks.readonly",
		},
		Endpoint: google.Endpoint,
	}

	client, err := getClient(oauthConfig)
	if err != nil {
		log.Fatalf("Unable to get client: %v\n", err)
	}
}

func redirectCallbackUrl(url string) {
	fmt.Println("Redirecting you to the login page...")

	err := browser.OpenURL(url)
	if err != nil {
		log.Printf("Failed to open browser automatically: %v", err)
		fmt.Println("Please try again :(")
	}
}

func resolveToken(config *oauth2.Config) (*oauth2.Token, error) {
	token, err := loadToken("token.json")
	if err == nil {
		if token.Valid() {
			fmt.Println("✅ Loaded saved token, no need to login again.")
			return token, nil
		}

		tokenSrc := oauthConfig.TokenSource(context.Background(), token)
		newToken, err := tokenSrc.Token()

		if err == nil {
			fmt.Println("🔄 Token refreshed successfully.")
			saveToken("token.json", newToken)
			return newToken, nil
		}

	}

	tokenChan := make(chan *oauth2.Token)

	http.HandleFunc("/oauth2callback", makeCallbackHandler(config, tokenChan))

	go func() {
		port := "6769"
		fmt.Println("Server started at http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":6769", nil))
	}()

	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	redirectCallbackUrl(url)

	token = <-tokenChan
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Println("✅ Authentication successful!")

	return token, nil
}

func makeCallbackHandler(config *oauth2.Config, tokenChan chan *oauth2.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := oauthConfig.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Failed in exchange code", http.StatusInternalServerError)
			fmt.Println("❌Sorry, there was an issue exchanging authentication code. ")
			return
		}

		saveToken("token.json", token)
		tokenChan <- token

		fmt.Fprint(w, "✅ Authentication Successful! You can close this window.")
	}
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}

func loadToken(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)

	return &token, err
}

func getClient(config *oauth2.Config) (*http.Client, error) {
	token, err := resolveToken(config)
	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), token), nil
}
