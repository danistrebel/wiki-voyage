package main

import (
	"testing"
)

func TestActivitityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		want     string
	}{
		{
			name:     "See activity",
			activity: "see",
			want:     "👀",
		},
		{
			name:     "Do activity",
			activity: "do",
			want:     "🤸",
		},
		{
			name:     "Eat activity",
			activity: "eat",
			want:     "🍽️",
		},
		{
			name:     "Sleep activity",
			activity: "sleep",
			want:     "😴",
		},
		{
			name:     "Buy activity",
			activity: "buy",
			want:     "🛍️",
		},
		{
			name:     "Drink activity",
			activity: "drink",
			want:     "🍹",
		},
		{
			name:     "Unknown activity",
			activity: "unknown",
			want:     "📍",
		},
		{
			name:     "Empty activity",
			activity: "",
			want:     "📍",
		},
		{
			name:     "Activity with different case",
			activity: "SEE",
			want:     "📍",
		},
		{
			name:     "Activity with spaces",
			activity: " see ",
			want:     "📍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := activitityEmoji(tt.activity); got != tt.want {
				t.Errorf("activitityEmoji(%q) = %q, want %q", tt.activity, got, tt.want)
			}
		})
	}
}
