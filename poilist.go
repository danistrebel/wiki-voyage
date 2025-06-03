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

// activityEmoji returns an emoji string for a given activity type.
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
		return "🍹"
	default:
		return "✨" // Default emoji for other activities
	}
}

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("Missing PROJECT_ID")
	}

	ctx := context.Background()

	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)

	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("error creating BigQuery client: %w", err)
	}
	defer client.Close()

	queryJob := client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error creating query job: %w", err)
	}

	var pointsOfInterest []PointOfInterest
	it, err := queryJob.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading query results: %w", err)
	}
	for {
		var row PointOfInterest
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}
		row.Icon = activityEmoji(row.Activity)
		pointsOfInterest = append(pointsOfInterest, row)
	}
	return pointsOfInterest, nil

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
