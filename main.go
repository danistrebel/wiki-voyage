package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", listPointsOfInterestHandler)
	http.HandleFunc("/recommendations", recommendationHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	chromeUserProfile := os.Getenv("GOOGLE_CHROME_USER_PROFILE")
	if chromeUserProfile == "" {
		chromeUserProfile = "1"
	}
	log.Printf("Running app on http://0.0.0.0:%s?authuser=%s", port, chromeUserProfile)
	http.ListenAndServe(":"+port, nil)
}
