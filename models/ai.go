package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type RequestBody struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

func GgtAI() error {
	openAIEndpoint := "https://api.openai.com/v1/completions"
	err := godotenv.Load()
	if err != nil {
		return err
	}

	apiKey := os.Getenv("GPT_KEY")

	requestBody := &RequestBody{
		Model:     "text-davinci-002",
		Prompt:    "Once upon a time,",
		MaxTokens: 50,
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(jsonValue))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var responseMap map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseMap)

	completion := responseMap["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
	println("Completion:", completion)
	return nil
}
