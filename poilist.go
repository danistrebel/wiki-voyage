package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/bigquery"
)

var (
	poiCache   = make(map[string][]PointOfInterest)
	poiCacheMu sync.RWMutex
)

// limitStringLength limits the length of a string to a specified maximum length.
// If the string's length exceeds the maximum length, it is truncated,
// and ".." is appended to indicate that the string has been shortened.
func limitStringLength(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + ".."
}

// activityEmoji maps activity types to emojis.
func activityEmoji(activity string) string {
	switch activity {
	case "see":
		return "👀"
	case "do":
		return "🤸"
	case "eat":
		return "🍽️"
	case "sleep":
		return "🛌"
	case "buy":
		return "🛍️"
	case "drink":
		return "🍹"
	default:
		return "🤌"
	}
}

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {

	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("Missing PROJECT_ID")
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)

	q := client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %w", err)
	}
	var pois []PointOfInterest
	for {
		var poi PointOfInterest
		err := it.Next(&poi)
		if err != nil {
			break
		}
		// Add emoji to the poi
		poi.Icon = activityEmoji(poi.Activity)
		poi.Title = limitStringLength(poi.Title, 30)              // Example: Limit Title to 30 characters
		poi.Description = limitStringLength(poi.Description, 100) // Example: Limit Description to 100 characters

		pois = append(pois, poi)
	}
	return pois, nil
}

const defaultCity = "Rome"

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
