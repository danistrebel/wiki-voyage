package main

import "testing"

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		want     string
	}{
		{
			name:     "See activity",
			activity: "see",
			want:     "🏞️",
		},
		{
			name:     "Do activity",
			activity: "do",
			want:     "🤸",
		},
		{
			name:     "Eat activity",
			activity: "eat",
			want:     "🍔",
		},
		{
			name:     "Drink activity",
			activity: "drink",
			want:     "🍹",
		},
		{
			name:     "Buy activity",
			activity: "buy",
			want:     "🛍️",
		},
		{
			name:     "Unknown activity",
			activity: "unknown",
			want:     "❓",
		},
		{
			name:     "Empty activity",
			activity: "",
			want:     "❓",
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
