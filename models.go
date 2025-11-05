package main

type PageData struct {
	PageTitle        string
	MapsApiKey       string
	PointsOfInterest []PointOfInterest
}

type PointOfInterest struct {
	Id          string // UUID4
	Title       string
	Latitude    float64
	Longitude   float64
	Description string
	Activity    string // values like "see", "do", "eat", "sleep", "buy", "drink"
	Icon        string
}

type RecommendationRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type RecommendationResponse struct {
	Recommendation string `json:"recommendation"`
}
