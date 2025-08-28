package main

import "testing"

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		want     string
	}{
		{"see", "see", "👀"},
		{"do", "do", "🤸"},
		{"eat", "eat", "🍽️"},
		{"sleep", "sleep", "😴"},
		{"buy", "buy", "🛍️"},
		{"drink", "drink", "🍻"},
		{"unknown", "unknown", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := activityEmoji(tt.activity); got != tt.want {
				t.Errorf("activityEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}
