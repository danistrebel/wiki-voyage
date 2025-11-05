package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"unicode/utf8"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var (
	poiCache   = make(map[string][]PointOfInterest)
	poiCacheMu sync.RWMutex
)

func limitStringLength(s string, length int) string {
	if utf8.RuneCountInString(s) > length {
		// Truncate and add ellipsis
		return string([]rune(s)[0:length-2]) + ".."
	}
	return s
}

func getEmojiForActivity(activity string) string {
	switch activity {
	case "see":
		return "👀"
	case "do":
		return "🤸"
	case "eat":
		return "🍔"
	case "sleep":
		return "😴"
	case "buy":
		return "🛍️"
	case "drink":
		return "🍹"
	default:
		return "📍"
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
		return nil, err
	}
	defer client.Close()
	query := fmt.Sprintf("SELECT * FROM `%s.wiki_voyage.points_of_interest` WHERE city = @city", projectId)
	q := client.Query(query)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "city",
			Value: city,
		},
	}
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}

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
		pois = append(pois, poi)
	}
	return pois, nil

}

const defaultCity = "Berlin"

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

	if len(PointsOfInterest) > 20 {
		PointsOfInterest = PointsOfInterest[0:20]
	}

	for i := range PointsOfInterest {
		PointsOfInterest[i].Description = limitStringLength(PointsOfInterest[i].Description, 100)
		PointsOfInterest[i].Icon = getEmojiForActivity(PointsOfInterest[i].Activity)
	}

	data := PageData{
		PageTitle:        fmt.Sprintf("Things to do in %s:", city),
		MapsApiKey:       MapsApiKey,
		PointsOfInterest: PointsOfInterest,
	}

	tmpl.Execute(w, data)
}
