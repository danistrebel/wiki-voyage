package main

import "testing"

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
			name:     "unknown activity",
			activity: "fly",
			want:     "",
		},
		{
			name:     "empty activity",
			activity: "",
			want:     "",
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

func TestLimitStringLength(t *testing.T) {
	tests := []struct {
		name      string
		s         string
		maxLength int
		want      string
	}{
		{"string shorter than max", "hello", 10, "hello"},
		{"string equal to max", "hello world", 11, "hello world"},
		{"string longer than max", "hello beautiful world", 10, "hello beau..."},
		{"empty string", "", 10, ""},
		{"zero max length", "hello", 0, "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := limitStringLength(tt.s, tt.maxLength); got != tt.want {
				t.Errorf("limitStringLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
