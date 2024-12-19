package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/sessions"
)

var fileService = &service.FileService{}
var aiService = &service.AIService{
	Client: &http.Client{},
}

var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

// HandleData handles data retrieval
func HandleData(w http.ResponseWriter, r *http.Request) {
	data, err := fileService.ReadCSV("data-series.csv")
	if err != nil {
		http.Error(w, "Failed to read data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func HandleAddData(w http.ResponseWriter, r *http.Request) {
	reqBody := service.EnergyData{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = fileService.AddDataCSV("data-series.csv", reqBody)

	if err != nil {
		http.Error(w, "Failed to add data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Data added successfully"})
}

// HandleAnalyticsEnergy handles analytics energy requests
func HandleAnalyticsEnergy(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Query string `json:"query"`
	}

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	}

	// Baca data dari file csv
	data, err := fileService.ReadCSV("data-series.csv")
	if err != nil {
		http.Error(w, "Failed to read data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Lakukan analisis data
	dailyConsumption, weeklyConsumption, monthlyConsumption := fileService.AnalyzeEnergyConsumption(data)
	roomConsumption := fileService.CalculateRoomEnergyConsumption(data)

	// Estimasi biaya listrik
	costPerKWh := 1500.0

	estimatedCost := fileService.EstimateElectricityCost(data, costPerKWh)

	response := map[string]interface{}{
		"daily_consumption":   dailyConsumption,
		"weekly_consumption":  weeklyConsumption,
		"monthly_consumption": monthlyConsumption,
		"room_consumption":    roomConsumption,
		"estimated_cost":      estimatedCost,
	}

	// Mengirimkan respons dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

// HandleAnalyzeAI handles AI-based analysis
func HandleAnalyzeAI(w http.ResponseWriter, r *http.Request, token string) {
	var request struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call AnalyzeData from AIService
	result, err := aiService.AnalyzeData("data-series.csv", request.Query, token)
	if err != nil {
		http.Error(w, "Failed to analyze data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status": "success",
		"result": result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleChat handles chat requests
func HandleChat(w http.ResponseWriter, r *http.Request, token string) {
	var request struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get chat session context
	session := getSession(r)
	context := ""
	if session.Values["context"] != nil {
		if value, ok := session.Values["context"].(string); ok {
			context = value
		} else {
			http.Error(w, "Invalid session context", http.StatusInternalServerError)
			log.Println("Invalid type for context; resetting to empty string.")
		}
	}

	// Communicate with the chat AI
	response, err := aiService.ChatWithAI(context, request.Query, token)
	if err != nil {
		http.Error(w, "Failed to chat with AI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update chat session context
	session.Values["context"] = context + "\n" + response.GeneratedText
	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to save session: %v\n", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Respond with the AI's response
	result := map[string]string{
		"status": "success",
		"answer": response.Answer,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func HandleQnA(w http.ResponseWriter, r *http.Request, token string) {
	var request struct {
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	answer, err := aiService.AnswerQuestion("data-series.csv", request.Question, token)
	if err != nil {
		http.Error(w, "Failed to qna with AI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}
