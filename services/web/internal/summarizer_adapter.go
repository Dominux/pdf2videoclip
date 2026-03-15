package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SummarizerAdapter struct{}

func newSummarizerAdapter() SummarizerAdapter {
	return SummarizerAdapter{}
}

type summarizerResponse struct {
	Response string `json:"response"`
}

func (s *SummarizerAdapter) generate(text string) string {
	data := map[string]any{"model": "qwen-zoomer", "prompt": text, "stream": false}
	jsonData, _ := json.Marshal(data)

	logger := log.Default()
	logger.Println("Before Summarizer request")

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Println("After Summarizer request")

	// 2. Check if the server actually sent a json (200 OK)
	if resp.StatusCode != http.StatusOK {
		panic("failed to download json: " + resp.Status)
	}

	// 2. Read all bytes from the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response summarizerResponse

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Response)

	return response.Response
}
