package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type HTTPBinResponse struct {
    Args    map[string]string `json:"args"`
    Headers map[string]string `json:"headers"`
    Origin  string            `json:"origin"`
    URL     string            `json:"url"`
}

func main() {
	fmt.Println("Starting program...")
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./http-get <url>\n")
		os.Exit(1)
	}

	targetURL := args[1]
	fmt.Printf("URL to fetch: %s\n", targetURL)

	if _, err := url.ParseRequestURI(targetURL); err != nil {
		fmt.Printf("Error: Invalid URL format - %s\n", err)
		os.Exit(1)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new request with context
	fmt.Println("Creating HTTP request...")
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	fmt.Println("Request created")

	// Make the HTTP request
	fmt.Println("Creating HTTP client...")
	client := &http.Client{
		Timeout: 10 * time.Second, // Total timeout for the entire request
	}

	fmt.Println("Sending HTTP request...")
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error in HTTP request: %v\n", err)
		log.Fatal(err)
	}
	fmt.Printf("Received response with status: %s\n", response.Status)

	defer response.Body.Close()

	fmt.Println("Reading response body...")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		log.Fatalf("Error reading response: %v", err)
	}
	fmt.Println("Response body read successfully")

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("Error: Received status code %d\n", response.StatusCode)
		log.Fatalf("Invalid status code %d: %s", response.StatusCode, string(body))
	}
	fmt.Println("Status code is valid (2xx)")

	var responseData HTTPBinResponse

	fmt.Println("Parsing JSON response...")
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		log.Fatal(err)
	}
	fmt.Println("JSON parsed successfully")

	// Print the response in a formatted way
	fmt.Println("\n=== Response ===")
	fmt.Printf("URL: %s\n", responseData.URL)
	fmt.Printf("Origin: %s\n", responseData.Origin)
	
	fmt.Println("\nHeaders:")
	for k, v := range responseData.Headers {
		fmt.Printf("  %s: %s\n", k, v)
	}

	if len(responseData.Args) > 0 {
		fmt.Println("\nArguments:")
		for k, v := range responseData.Args {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
}