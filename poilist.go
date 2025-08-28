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

func limitStringLength(s string, maxLength int) string {
	runes := []rune(s)
	if len(runes) <= maxLength {
		return s
	}
	return string(runes[:maxLength]) + ".."
}

func loadPointsOfInterest(ctx context.Context, city string) ([]PointOfInterest, error) {
	projectId := os.Getenv("PROJECT_ID")
	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = '%s'", projectId, city)
	it, _ := client.Query(query).Read(ctx)
	var pois []PointOfInterest
	for {
		var poi PointOfInterest
		err := it.Next(&poi)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		poi.Description = limitStringLength(poi.Description, 160)
		poi.Icon = activityEmoji(poi.Activity)
		pois = append(pois, poi)
	}
	return pois, nil

}

const defaultCity = "Uppsala"

func listPointsOfInterestHandler(w http.ResponseWriter, r *http.Request) {
	MapsApiKey := os.Getenv("MAPS_KEY")
	if MapsApiKey == "" {
		log.Fatal("Missing MAPS_KEY")
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
