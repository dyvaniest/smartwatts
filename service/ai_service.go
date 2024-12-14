package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(csvFilePath, query, token string) (string, error) {
	// Membaca data dari file CSV
	table, err := s.readCSVAsTable(csvFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(table) == 0 {
		return "", fmt.Errorf("table cannot be empty")
	}

	requestBody := map[string]interface{}{
		"inputs": map[string]interface{}{
			"query": query,
			"table": table,
		},
	}

	reqBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI model returned non-200 status: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read AI model response: %w", err)
	}
	log.Printf("AI Model Response: %s", string(body))

	var result model.TapasResponse
	if json.Unmarshal(body, &result) == nil && len(result.Cells) > 0 {
		return result.Cells[0], nil
	}

	return result.Cells[0], nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	if s == nil {
		log.Println("AIService is nil!")
		return model.ChatResponse{}, fmt.Errorf("AIService is not initialized")
	}

	// Cek token
	if token == "" {
		log.Println("Hugging Face token is empty!")
		return model.ChatResponse{}, fmt.Errorf("hugging Face token is not provided")
	}

	// Gunakan hanya query sebagai input
	data := map[string]interface{}{
		"inputs": query,
		"parameters": map[string]interface{}{
			"max_new_tokens": 300,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal request data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return model.ChatResponse{}, fmt.Errorf("AI model returned non-200 status: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to read AI model response: %w", err)
	}

	log.Printf("AI Model Response: %s", string(body))

	var result []model.ChatResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to parse AI model response: %w", err)
	}

	if len(result) == 0 || result[0].GeneratedText == "" {
		return model.ChatResponse{}, fmt.Errorf("AI model response is empty or invalid")
	}

	// Kembalikan hanya teks yang dihasilkan oleh AI
	return model.ChatResponse{Answer: result[0].GeneratedText}, nil
}

func (s *AIService) readCSVAsTable(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("insufficient data in CSV")
	}

	headers := records[0]
	table := make(map[string][]string)

	for _, header := range headers {
		table[header] = []string{}
	}

	for _, row := range records[1:] {
		for i, value := range row {
			table[headers[i]] = append(table[headers[i]], value)
		}
	}

	return table, nil
}
