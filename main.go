package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var question string

func init() {
	flag.StringVar(&question, "ques", "", "question to ask")
	flag.Parse()
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get API key from environment
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found in environment")
	}

	// Define the API endpoint
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	// Create the JSON payload
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": question},
				},
			},
		},
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	// Create a POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	// Extract the text
	candidates, ok := responseBody["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		log.Fatal("No candidates found in response")
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		log.Fatal("No content found in first candidate")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		log.Fatal("No parts found in content")
	}

	text, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		log.Fatal("No text found in parts")
	}

	// Print the extracted text
	fmt.Println("Response:", text)
}
