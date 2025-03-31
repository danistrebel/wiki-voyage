package main

import (
	"testing"
)

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		want     string
	}{
		{
			name:     "see activity",
			activity: "see",
			want:     "👀",
		},
		{
			name:     "do activity",
			activity: "do",
			want:     "🤸",
		},
		{
			name:     "eat activity",
			activity: "eat",
			want:     "🍽️",
		},
		{
			name:     "sleep activity",
			activity: "sleep",
			want:     "🛌",
		},
		{
			name:     "buy activity",
			activity: "buy",
			want:     "🛍️",
		},
		{
			name:     "drink activity",
			activity: "drink",
			want:     "🍹",
		},
		{
			name:     "unknown activity",
			activity: "unknown",
			want:     "📍",
		},
		{
			name:     "empty activity",
			activity: "",
			want:     "📍",
		},
		{
			name:     "mixed case activity",
			activity: "SeE",
			want:     "📍",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := activityEmoji(tt.activity); got != tt.want {
				t.Errorf("activityEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}
