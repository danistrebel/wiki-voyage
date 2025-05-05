package main

import "testing"

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		activity string
		expected string
	}{
		{"see", "👀"},
		{"do", "🤸‍♀️"},
		{"eat", "🍽️"},
		{"sleep", "🛌"},
		{"buy", "🛍️"},
		{"drink", "🍹"},
		{"unknown", ""},
		{"", ""}, // Test empty string
	}

	for _, test := range tests {
		actual := activityEmoji(test.activity)
		if actual != test.expected {
			t.Errorf("For activity '%s', expected emoji '%s', but got '%s'", test.activity, test.expected, actual)
		}
	}
}
