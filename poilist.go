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

func loadPointsOfInterest(city string) ([]PointOfInterest, error) {
	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("Missing PROJECT_ID")
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectId)
	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	queryJob := client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error creating BigQuery query: %w", err)
	}

	var pointsOfInterest []PointOfInterest
	it, err := queryJob.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading BigQuery query result: %w", err)
	}
	for {
		var row PointOfInterest
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating BigQuery result: %w", err)
		}
		row.Icon = activityEmoji(row.Activity)

		row.Description = limitStringLength(row.Description, 100)

		pointsOfInterest = append(pointsOfInterest, row)
	}
	return pointsOfInterest, nil
}

// function to represent activity types as emojis
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

// Calculate the approximate distance between two places based on geo coordinates

// limitStringLength truncates a string to a maximum length.
func limitStringLength(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	// Truncate and add an ellipsis.
	return s[:maxLength] + "..."
}

const defaultCity = "Munich"

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

	// Truncate descriptions before rendering
	for i := range PointsOfInterest {
		PointsOfInterest[i].Description = limitStringLength(PointsOfInterest[i].Description, 100)
	}

	data := PageData{
		PageTitle:        fmt.Sprintf("Things to do in %s:", city),
		MapsApiKey:       MapsApiKey,
		PointsOfInterest: PointsOfInterest,
	}

	tmpl.Execute(w, data)
}
