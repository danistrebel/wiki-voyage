package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitStringLength(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"Short string", 20, "Short string"},
		{"Long string that needs to be truncated", 10, "Long strin.."},
		{"Exact length string", 19, "Exact length string"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual := limitStringLength(test.input, test.maxLength)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestActivityEmoji(t *testing.T) {
	assert.Equal(t, "👀", activityEmoji("see"))
	assert.Equal(t, "🤸", activityEmoji("do"))
	assert.Equal(t, "🍽️", activityEmoji("eat"))
	assert.Equal(t, "😴", activityEmoji("sleep"))
	assert.Equal(t, "🛍️", activityEmoji("buy"))
	assert.Equal(t, "🍻", activityEmoji("drink"))
	assert.Equal(t, "", activityEmoji("unknown"))
}
