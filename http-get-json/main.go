package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"context"
)

type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	// Create a context with timeout

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new request with context
    req, err:= http.NewRequestWithContext(ctx, "GET", args[1], nil)
    if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Make the HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second, // Total timeout for the entire request
	}

	response, err := client.Do(req)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatal("Request timed out")
		}
		log.Fatal("Request failed %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatalf("Invalid status code %d: %s", response.StatusCode, string(body))
	}
    

	var words Words

	err = json.Unmarshal(body, &words)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("JSON: Parsed:\nPage: %s\nWords: %s\n", words.Page, strings.Join(words.Words, ", "))
}