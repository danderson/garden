package views

import (
	"testing"
	"time"

	"go.universe.tf/garden/types"
)

func TestDaysBetween(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"2024-01-21 10:40:00", "2024-01-21 10:37:45", 0},
		{"2024-01-20 10:40:00", "2024-01-21 10:37:45", 1},
		{"2024-01-20 23:40:00", "2024-01-21 01:37:45", 1},
		{"2024-01-15 23:40:00", "2024-01-21 01:37:45", 6},
		{"2023-12-30 12:28:23", "2024-01-02 01:37:45", 3},
		{"2024-01-19 10:30:00", "2024-01-20 00:00:00", 1},
		{"2024-01-19 00:00:00", "2024-01-20 00:00:00", 1},
		{"2024-01-20 00:00:00", "2024-01-20 00:00:00", 0},
	}

	for _, tc := range cases {
		a, err := time.ParseInLocation("2006-01-02 15:04:05", tc.a, types.Pacific)
		if err != nil {
			t.Fatal(err)
		}
		b, err := time.ParseInLocation("2006-01-02 15:04:05", tc.b, types.Pacific)
		if err != nil {
			t.Fatal(err)
		}
		got := daysBetween(a, b)
		if got != tc.want {
			t.Errorf("daysBetween(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
		}
		got2 := daysBetween(b, a)
		if got != got2 {
			t.Errorf("daysBetween(%q, %q) != daysBetween(%q, %q)", tc.a, tc.b, tc.b, tc.a)
		}
	}
}
