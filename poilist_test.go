package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivityEmoji(t *testing.T) {
	tests := []struct {
		activity string
		want     string
	}{
		{"see", "👀"},
		{"do", "🤸"},
		{"eat", "🍽️"},
		{"sleep", "😴"},
		{"buy", "🛍️"},
		{"drink", "🍻"},
		{"other", ""},
	}

	for _, tt := range tests {
		t.Run(tt.activity, func(t *testing.T) {
			got := activityEmoji(tt.activity)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestActivityEmoji_Empty(t *testing.T) {
	got := activityEmoji("")
	assert.Equal(t, "", got)
}

func TestActivityEmoji_Unknown(t *testing.T) {
	got := activityEmoji("unknown")
	assert.Equal(t, "", got)
}
