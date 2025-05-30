package main

import "testing"

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		activity string
		expected string
	}{
		{"see", "👀"},
		{"do", "🤸"},
		{"eat", "🍔"},
		{"sleep", "😴"},
		{"buy", "🛍️"},
		{"drink", "🍹"},
		{"unknown", "📍"},          // Test default case
		{"", "📍"},                 // Test empty string
		{"SEE", "👀"},              // Test uppercase
		{"  see  ", "👀"},          // Test leading/trailing spaces
		{"  UNKNOWN  ", "📍"},      // Test uppercase and spaces with default
	}

	for _, test := range tests {
		t.Run(test.activity, func(t *testing.T) {
			result := activityEmoji(test.activity)
			if result != test.expected {
				t.Errorf("activityEmoji(%s) = %s; want %s", test.activity, result, test.expected)
			}
		})
	}
}

func TestActivityEmoji_DefaultCase(t *testing.T) {
	expected := "📍"
	result := activityEmoji("non_existent_activity")
	if result != expected {
		t.Errorf("activityEmoji with default case failed: got %s, want %s", result, expected)
	}
}

func TestActivityEmoji_EmptyString(t *testing.T) {
	expected := "📍"
	result := activityEmoji("")
	if result != expected {
		t.Errorf("activityEmoji with empty string failed: got %s, want %s", result, expected)
	}
}
