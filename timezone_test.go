package timezoneLookup

import (
	"os"
	"testing"
)

var (
	searchTestCases = map[[2]float64]string{
		// Examples from https://github.com/ringsaturn/tzf
		{34.4200, 111.8674}: "Asia/Shanghai",
		{34.4200, -97.8674}: "America/Chicago",
		{31.1139, 121.3547}: "Asia/Shanghai",
		{36.4432, 139.4382}: "Asia/Tokyo",
		{50.2506, 24.5212}:  "Europe/Kiev",
		{52.0152, -0.9671}:  "Europe/London",
		{46.2747, -4.5706}:  "Etc/GMT",
		{45.0182, 111.9781}: "Asia/Shanghai",
		{38.3530, -73.7729}: "Etc/GMT+5",

		{37.7749, -122.4194}: "America/Los_Angeles",
	}
)

func buildCache() (Timezonecache, error) {
	var tzc Timezonecache
	f, err := os.Open("timezone.data")
	if err != nil {
		return tzc, err
	}
	return tzc, tzc.Load(f)
}

func TestSearch(t *testing.T) {
	tzc, err := buildCache()
	if err != nil {
		t.Error(err)
	}

	for point, name := range searchTestCases {
		result, err := tzc.Search(point[0], point[1])
		if err != nil {
			t.Error(err)
		}
		if result.Name != name {
			t.Errorf("Expected: %v, got %v", name, result.Name)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	tzc, err := buildCache()
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for point := range searchTestCases {
			_, err := tzc.Search(point[0], point[1])
			if err != nil {
				b.Error(err)
			}
		}
	}
}
