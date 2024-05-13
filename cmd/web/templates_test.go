package main

import (
	"testing"
	"time"

	"github.com/rynhndrcksn/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// Use an anonymous struct to create all of our test cases.
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 05, 13, 11, 0, 0, 0, time.UTC),
			want: "13 May 2024 at 11:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 05, 13, 11, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "13 May 2024 at 10:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
