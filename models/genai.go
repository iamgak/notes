package models

import (
	"context"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	API_KEY       = ""  // We'll go back to this later
	batchWindowMs = 100 // add whatever value you want, this will wait for 100ms
)

var systemMessage = openai.ChatCompletionMessage{
	Role:    "system",
	Content: `You are a helpful assistant. Transform the query into a boolean logic using boolean operators (AND, OR, NOT...), don't show comments.`,
}

type ResponseBuffer struct {
	response string
}

type RequestBuffer struct {
	client        *openai.Client
	batchWindow   time.Duration // How long the buffer will wait before sending the request
	buffer        []bufferedRequest
	processing    bool
	lastRequestAt time.Time
	mu            sync.Mutex // We use a mutex ensure only what is supposed to access this buffer has access to it
}

type bufferedRequest struct {
	query    string
	resultCh chan string
}

type QueryRequest struct {
	Parameter string `json:"parameter"`
}

// Setting up the requestBuffer

func NewRequestBuffer(client *openai.Client) *RequestBuffer {
	return &RequestBuffer{
		client:      client,
		batchWindow: time.Duration(batchWindowMs) * time.Millisecond,
		buffer:      make([]bufferedRequest, 0), // Creates an empty buffer
	}
}

func (rb *RequestBuffer) AddRequest(query string) (string, error) {
	rb.mu.Lock()
	currentTime := time.Now()
	resultCh := make(chan string, 1) // Creates a new channel for the response

	rb.buffer = append(rb.buffer, bufferedRequest{
		query:    query,
		resultCh: resultCh,
	})

	// If the 100ms have passed, it processes the batch (or single request)
	if currentTime.Sub(rb.lastRequestAt) > rb.batchWindow && !rb.processing {
		rb.processing = true
		go rb.processBatch()
	}

	rb.lastRequestAt = currentTime
	rb.mu.Unlock() // Unlocks the buffer so the next routine can access it

	result := <-resultCh
	return result, nil
}

func (rb *RequestBuffer) processBatch() {
	defer func() {
		rb.mu.Lock()
		rb.processing = false
		rb.mu.Unlock()
	}()

	time.Sleep(rb.batchWindow)

	rb.mu.Lock()
	if len(rb.buffer) == 0 {
		rb.mu.Unlock()
		return
	}

	requests := make([]bufferedRequest, len(rb.buffer))
	copy(requests, rb.buffer)              // Copy buffer to ensure thread safety
	rb.buffer = make([]bufferedRequest, 0) // Clear the buffer
	rb.mu.Unlock()

	ctx := context.Background()

	if len(requests) == 1 {
		// Single request processing
		resp, err := rb.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: "gpt-4o-mini",
				Messages: []openai.ChatCompletionMessage{
					systemMessage,
					{
						Role:    "user",
						Content: requests[0].query,
					},
				},
			},
		)

		if err != nil {
			requests[0].resultCh <- "Error processing request"
			return
		}

		requests[0].resultCh <- resp.Choices[0].Message.Content
	} else {
		// Batch processing
		resp, err := rb.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: "gpt-4o-mini",
				Messages: []openai.ChatCompletionMessage{
					systemMessage,
					{
						Role:    "user",
						Content: requests[0].query,
					},
				},
				N: len(requests),
			},
		)

		if err != nil {
			for _, req := range requests {
				req.resultCh <- "Error processing batch request"
			}
			return
		}

		for i, choice := range resp.Choices {
			if i < len(requests) {
				requests[i].resultCh <- choice.Message.Content
			}
		}
	}
}

func init() {

}

func Prompt(req string) (ResponseBuffer, error) {
	client := openai.NewClient(API_KEY)
	requestBuffer := NewRequestBuffer(client)
	// Spins up a new openaiclient, verified with our api_key

	result, err := requestBuffer.AddRequest(req)
	resp := ResponseBuffer{
		response: result,
	}

	return resp, err
}
