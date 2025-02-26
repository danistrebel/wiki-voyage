package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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

	generatedRecommendation := fmt.Sprintf("I found this place called %s. You absolutely need to check it out!", req.Title)

	resp := RecommendationResponse{Recommendation: generatedRecommendation}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println("Error encoding response:", err)
		return
	}
}
