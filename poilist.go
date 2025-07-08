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

func limitStringLength(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {
	project := os.Getenv("PROJECT_ID")
	if project == "" {
		log.Fatal("Missing PROJECT_ID")
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	q := client.Query("SELECT * FROM `wiki_voyage.points_of_interest` WHERE city = @city")
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "city",
			Value: city,
		},
	}

	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var pois []PointOfInterest
	for {
		var poi PointOfInterest
		err := it.Next(&poi)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}
		poi.Icon = activityEmoji(poi.Activity)
		poi.Description = limitStringLength(poi.Description, 80)
		pois = append(pois, poi)
	}
	return pois, nil
}

// activityEmoji returns an emoji for a given activity type.
func activityEmoji(activity string) string {
	switch activity {
	case "see":
		return "👀"
	case "do":
		return "🤸"
	case "eat":
		return "🍽️"
	case "sleep":
		return "😴"
	case "buy":
		return "🛍️"
	case "drink":
		return "🍻"
	default:
		return ""
	}
}

const defaultCity = "London"

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
