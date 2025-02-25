package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

// Use Gemini Flash 2.0 to create recommendation about about a place to visit
func CreatePlaceRecommendation(title, description string) (string, error) {
	projectId := os.Getenv("PROJECT_ID")

	location := "europe-west1"
	modelName := "gemini-2.0-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectId, location)
	if err != nil {
		return "", fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)
	prompt := genai.Text(fmt.Sprintf(`
		You are a well-travelled individual. You just remembered a place that you want to share with a friend and suggest them to visit it.
		Write a short recommendation for the following place. Only write the message itself and skip any greetings:

		Title %s
		Description: %s
		`, title, description))

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(resp.Candidates[0].Content.Parts[0]), nil
	} else {
		return "", fmt.Errorf("no content generated")
	}
}

func recommendationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RecommendationRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Println("Error reading request body:", err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		return

	}

	generatedRecommendation, _ := CreatePlaceRecommendation(req.Title, req.Description)

	resp := RecommendationResponse{Recommendation: generatedRecommendation}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}
}
