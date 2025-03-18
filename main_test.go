package main

import (
	"testing"
)

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		expected string
	}{
		{
			name:     "See Activity",
			activity: "see",
			expected: "👀",
		},
		{
			name:     "Do Activity",
			activity: "do",
			expected: "🤸",
		},
		{
			name:     "Eat Activity",
			activity: "eat",
			expected: "🍽️",
		},
		{
			name:     "Sleep Activity",
			activity: "sleep",
			expected: "😴",
		},
		{
			name:     "Buy Activity",
			activity: "buy",
			expected: "🛍️",
		},
		{
			name:     "Drink Activity",
			activity: "drink",
			expected: "🍹",
		},
		{
			name:     "Unknown Activity",
			activity: "unknown",
			expected: "🤡",
		},
		{
			name:     "Empty Activity",
			activity: "",
			expected: "🤡",
		},
		{
			name:     "Numbers as Activity",
			activity: "123",
			expected: "🤡",
		},
		{
			name:     "Special chars as Activity",
			activity: "!@#$",
			expected: "🤡",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := activityEmoji(tt.activity)
			if actual != tt.expected {
				t.Errorf("activityEmoji(%s) = %s; expected %s", tt.activity, actual, tt.expected)
			}
		})
	}
}
