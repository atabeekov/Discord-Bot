package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const wikiRecentChangesAPI = "https://en.wikipedia.org/w/api.php?action=query&list=recentchanges&rclimit=10&rcprop=title|ids|sizes|flags|user&format=json"

// RecentChange represents a single recent change from the API
type RecentChange struct {
	Title string `json:"title"`
	User  string `json:"user"`
}

// WikiAPIResponse represents the structure of the response from Wikimedia API
type WikiAPIResponse struct {
	Query struct {
		RecentChanges []RecentChange `json:"recentchanges"`
	} `json:"query"`
}

// GetRecentChanges fetches recent changes from Wikimedia REST API
func GetRecentChanges() (string, error) {
	log.Println("Connecting to Wikipedia REST API...")

	resp, err := http.Get(wikiRecentChangesAPI)
	if err != nil {
		log.Println("Error connecting to API:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var apiResponse WikiAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	result := "Recent changes:\n"
	for _, change := range apiResponse.Query.RecentChanges {
		result += fmt.Sprintf("- Title: %s, User: %s\n", change.Title, change.User)
	}

	log.Println("Successfully fetched recent changes.")
	return result, nil
}

// PollRecentChanges polls Wikipedia REST API for recent changes
func PollRecentChanges(callback func(string)) {
	for {
		changes, err := GetRecentChanges()
		if err != nil {
			log.Println("Error fetching Wikipedia changes:", err)
		} else {
			callback(changes) // Pass data to the bot
		}
		time.Sleep(10 * time.Minute)
	}
}
