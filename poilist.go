package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	poiCache   = make(map[string][]PointOfInterest)
	poiCacheMu sync.RWMutex
)

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {

	return []PointOfInterest{
		{
			Title:       "Google ZRH Europaallee",
			Description: "The Google office in Zurich is an engineering hub for artificial intelligence, machine learning, and natural language processing. The office is also home to teams working on Google products such as Gemini, Maps, and YouTube.",
			Latitude:    47.3789437,
			Longitude:   8.5324559,
			Activity:    "work",
			Icon:        "🤓",
		},
		{
			Title:       "Google ZRH Brandschenkestrasse",
			Description: "Another Google office in Zurich.",
			Latitude:    47.365464,
			Longitude:   8.525309,
			Activity:    "work",
			Icon:        "🤓",
		},
	}, nil
}

const defaultCity = "Zurich"

func listPointsOfInterestHandler(w http.ResponseWriter, r *http.Request) {
	MapsApiKey := os.Getenv("MAPS_KEY")
	if MapsApiKey == "" {
		log.Fatal("Missing MAPS_API_KEY")
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	city := r.URL.Query().Get("city")
	if city == "" {
		city = defaultCity
	}

	poiCacheMu.RLock()
	PointsOfInterest, ok := poiCache[city]
	poiCacheMu.RUnlock()

	if !ok {
		result, err := loadPointsOfInterest(city)
		if err != nil {
			log.Println("Error loading points of interest:", err)
			http.Error(w, "Oops, there was an error loading the points of interest", http.StatusInternalServerError)
			return
		}
		PointsOfInterest = result
		poiCacheMu.Lock()
		poiCache[city] = PointsOfInterest
		poiCacheMu.Unlock()
	}

	data := PageData{
		PageTitle:        fmt.Sprintf("Things to do in %s:", city),
		MapsApiKey:       MapsApiKey,
		PointsOfInterest: PointsOfInterest,
	}

	tmpl.Execute(w, data)
}
