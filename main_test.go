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
			name:     "see",
			activity: "see",
			want:     "👀",
		},
		{
			name:     "do",
			activity: "do",
			want:     "🤸",
		},
		{
			name:     "eat",
			activity: "eat",
			want:     "🍽️",
		},
		{
			name:     "sleep",
			activity: "sleep",
			want:     "😴",
		},
		{
			name:     "buy",
			activity: "buy",
			want:     "🛍️",
		},
		{
			name:     "drink",
			activity: "drink",
			want:     "🍹",
		},
		{
			name:     "unknown",
			activity: "unknown",
			want:     "📍",
		},
		{
			name:     "empty",
			activity: "",
			want:     "📍",
		},
		{
			name:     "uppercase",
			activity: "SEE",
			want:     "📍",
		},
		{
			name:     "number",
			activity: "123",
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
