package main

import (
	"testing"
)

func TestLimitStringLength(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		length int
		want   string
	}{
		{
			name:   "String shorter than limit",
			s:      "Hello",
			length: 10,
			want:   "Hello",
		},
		{
			name:   "String equal to limit",
			s:      "Hello World",
			length: 11,
			want:   "Hello World",
		},
		{
			name:   "String longer than limit",
			s:      "This is a long string",
			length: 10,
			want:   "This is ..",
		},
		{
			name:   "String with unicode characters",
			s:      "Hello 👋 World",
			length: 10,
			want:   "Hello 👋 ..",
		},
		{
			name:   "Empty string",
			s:      "",
			length: 10,
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := limitStringLength(tt.s, tt.length); got != tt.want {
				t.Errorf("limitStringLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
