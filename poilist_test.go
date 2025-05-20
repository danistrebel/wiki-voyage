package main

import "testing"

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		activity string
		expected string
	}{
		{
			name:     "see",
			activity: "see",
			expected: "👀",
		},
		{
			name:     "do",
			activity: "do",
			expected: "🤸‍♀️",
		},
		{
			name:     "eat",
			activity: "eat",
			expected: "🍕",
		},
		{
			name:     "sleep",
			activity: "sleep",
			expected: "🛌",
		},
		{
			name:     "buy",
			activity: "buy",
			expected: "🛍️",
		},
		{
			name:     "drink",
			activity: "drink",
			expected: "🍻",
		},
		{
			name:     "unknown",
			activity: "unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := activityEmoji(tt.activity); got != tt.expected {
				t.Errorf("activityEmoji(%q) = %q, want %q", tt.activity, got, tt.expected)
			}
		}) 
	}
}

func TestActivityEmojiEdgeCases(t *testing.T) {
	if got := activityEmoji(""); got != "" {
		t.Errorf("activityEmoji(\"\") = %q, want \"\"", got)
	}
	if got := activityEmoji("  "); got != "" {
		t.Errorf("activityEmoji(\"  \") = %q, want \"\"", got)
	}
}
