package metrics

import (
	"testing"

	"github.com/avp-cloud/sermon/internal/pkg/models"
)

func TestGetStatus(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		upCodes string
	}{
		{
			name:    "succeeds for valid case",
			url:     "https://google.com",
			upCodes: "200,302",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status := GetStatus(tc.url, tc.upCodes)
			if status != models.StatusUp {
				t.Fatalf("GetStatus failed: got: %v want: UP", status)
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		upCodes string
	}{
		{
			name: "succeeds for valid case",
			url:  "https://google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GetMetrics(tc.url)
			if err != nil {
				t.Fatalf("GetMetrics err: got: %v want: nil", err)
			}
		})
	}
}
