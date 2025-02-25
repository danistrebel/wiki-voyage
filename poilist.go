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
	"google.golang.org/api/iterator"
)

var (
	poiCache   = make(map[string][]PointOfInterest)
	poiCacheMu sync.RWMutex
)

func activityEmoji(activity string) string {
	switch activity {
	case "see":
		return "👀" // Eyes
	case "do":
		return "🤸" // Person doing cartwheel
	case "eat":
		return "🍽️" // Fork and knife with plate
	case "sleep":
		return "🛌" // Person in bed
	case "buy":
		return "🛍️" // Shopping bags
	case "drink":
		return "🍹" // Tropical drink
	default:
		return "" // No emoji for unknown activity
	}
}

func limitStringLength(s string, maxLength int) (string, error) {
	if maxLength < 0 {
		return "", fmt.Errorf("maxLength cannot be negative")
	}
	if len(s) <= maxLength {
		return s, nil
	}
	return s[:maxLength] + "..", nil
}

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		return nil, fmt.Errorf("Missing PROJECT_ID")
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectId)
	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()
	q := client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query.Read: %w", err)
	}
	var pois []PointOfInterest
	for {
		var poi PointOfInterest
		err := it.Next(&poi)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterator.Next: %w", err)
		}
		poi.Icon = activityEmoji(poi.Activity)
		poi.Description, err = limitStringLength(poi.Description, 100)
		if err != nil {
			return nil, fmt.Errorf("limitStringLength: %w", err)
		}
		pois = append(pois, poi)
	}
	return pois, nil

}

const defaultCity = "Istanbul"

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
