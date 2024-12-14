package service_test

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var _ = Describe("AIService", func() {
	var (
		mockClient  *MockClient
		aiService   *service.AIService
		fileTestCSV string
	)

	BeforeEach(func() {
		mockClient = &MockClient{}
		aiService = &service.AIService{Client: mockClient}

		fileTestCSV = "test_data.csv"
		fileContent := "Date,Appliance,Energy_Consumption\n2024-12-10,EVCar,99.9\nRefrigerator,1.2\n"
		err := os.WriteFile(fileTestCSV, []byte(fileContent), 0644)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := os.Remove(fileTestCSV)
		Expect(err).To(BeNil())
	})

	Describe("AnalyzeData", func() {
		It("should analyze data from CSV and return a response", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				Expect(req.Method).To(Equal("POST"))
				Expect(req.URL.String()).To(ContainSubstring("https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq"))

				body, _ := io.ReadAll(req.Body)
				Expect(strings.Contains(string(body), "EVCar")).To(BeTrue())

				response := model.TapasResponse{
					Cells: []string{"99.9"},
				}
				respBody, _ := json.Marshal(response)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			}
			result, err := aiService.AnalyzeData(fileTestCSV, "What is the consumption for EVCar?", "mock-token")
			Expect(err).To(BeNil())
			Expect(result).To(Equal("99.9"))
		})

		It("should return an error if the API returns a non-200 status", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("Bad Request")),
				}, nil
			}

			result, err := aiService.AnalyzeData(fileTestCSV, "What is the consumption for EVCar?", "mock-token")
			Expect(err).To(BeNil())
			Expect(result).To(BeEmpty())
		})
	})

	Describe("ChatWithAI", func() {
		It("should return a chat response from the AI model", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				Expect(req.Method).To(Equal("POST"))
				Expect(req.URL.String()).To(ContainSubstring("https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct"))

				response := []model.ChatResponse{
					{GeneratedText: "The answer is 42."},
				}
				respBody, _ := json.Marshal(response)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBody)),
				}, nil
			}

			query := "What is the answer to life, the universe, and everything?"
			result, err := aiService.ChatWithAI("", query, "mock-token")
			Expect(err).To(BeNil())
			Expect(result.GeneratedText).To(Equal("The answer is 42."))
		})

		It("should return an error if the API response is invalid", func() {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("[]")),
				}, nil
			}

			result, err := aiService.ChatWithAI("", "Test query", "mock-token")
			Expect(err).To(BeNil())
			Expect(result.GeneratedText).To(BeEmpty())
		})
	})

})
