package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

var (
	poiCache   = make(map[string][]PointOfInterest)
	poiCacheMu sync.RWMutex
)

// limitStringLength truncates a string to a certain length and appends ".."
func limitStringLength(s string, length int) string {
	if len(s) > length {
		return s[:length] + ".."
	}
	return s
}

func getActivityEmoji(activity string) string {
	switch activity {
	case "see":
		return "🏞️"
	case "do":
		return "🤸"
	case "eat":
		return "😋"
	case "sleep":
		return "😴"
	case "buy":
		return "🛍️"
	case "drink":
		return "🍹"
	default:
		return "✨"
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
		return nil, fmt.Errorf("error creating bigquery client: %w", err)
	}
	defer client.Close()

	queryJob := client.Query("SELECT * FROM `" + projectId + ".wiki_voyage.points_of_interest` WHERE city = @city")
	queryJob.Parameters = []bigquery.QueryParameter{
		{
			Name:  "city",
			Value: city,
		},
	}
	it, err := queryJob.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %w", err)
	}

	var pointsOfInterest []PointOfInterest
	for {
		var row PointOfInterest
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating through results: %w", err)
		}
		row.Description = limitStringLength(row.Description, 140)
		row.Icon = getActivityEmoji(row.Activity)
		pointsOfInterest = append(pointsOfInterest, row)
	}
	return pointsOfInterest, nil

}

const defaultCity = "Utrecht"

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
