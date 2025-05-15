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
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + ".."
}

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

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("Missing PROJECT_ID")
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectId)
	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)
	if err != nil {
		log.Println("Error creating bigquery client:", err)
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()
	q := client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		log.Println("Error reading query results:", err)
		return nil, fmt.Errorf("query.Read: %v", err)
	}
	var pointsOfInterest []PointOfInterest
	for {
		var poi PointOfInterest
		err := it.Next(&poi)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Error iterating over query results:", err)
			return nil, fmt.Errorf("iteration error: %v", err)
		}
		poi.Icon = activityEmoji(poi.Activity)
		poi.Title = limitStringLength(poi.Title, 50)
		poi.Description = limitStringLength(poi.Description, 100)
		pointsOfInterest = append(pointsOfInterest, poi)
	}
	return pointsOfInterest, nil

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
