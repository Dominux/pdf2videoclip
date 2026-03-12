package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SummarizerAdapter struct{}

func newSummarizerAdapter() SummarizerAdapter {
	return SummarizerAdapter{}
}

func (s *SummarizerAdapter) generate(text string) {
	data := map[string]any{"model": "qwen-zoomer", "prompt": text, "stream": false}
	jsonData, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 2. Check if the server actually sent a file (200 OK)
	if resp.StatusCode != http.StatusOK {
		panic("failed to download file: " + resp.Status)
	}

	// 2. Read all bytes from the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 3. Convert bytes to string
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}
