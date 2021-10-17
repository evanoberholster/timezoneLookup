package timezoneLookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	WithSnappy = true
	NoSnappy   = false

	// Errors
	errNotExistGeoJSON    = "Error: GeoJSON file does not exist"
	errExistDatabase      = "Error: Destination Database file already exists"
	errNotExistDatabase   = "Error: Database file does not exist"
	errPolygonNotFound    = "Error: Polygon for Timezone not found"
	errTimezoneNotFound   = "Error: Timezone not found"
	errDatabaseTypeUknown = "Error: Database type unknown"
)

type TimezoneInterface interface {
	CreateTimezones(jsonFilename string) error
	LoadTimezones() error
	Query(q Coord) (string, error)
	Close()
}

type TimezoneGeoJSON struct {
	Type     string `json:"type"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Tzid string `json:"tzid"`
		} `json:"properties"`
		Geometry struct {
			Item        string        `json:"type"`
			Coordinates []interface{} `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

type Timezone struct {
	Tzid     string    `json:"tzid"`
	Polygons []Polygon `json:"polygons"`
}

type Polygon struct {
	Max    Coord   `json:"max"`
	Min    Coord   `json:"min"`
	Coords []Coord `json:"coords"`
}

type Coord struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Config struct {
	DatabaseName string
	DatabaseType string
	Snappy       bool
	Encoding     string
}

var Tz TimezoneInterface

func LoadTimezones(config Config) (TimezoneInterface, error) {
	if config.DatabaseType == "memory" {
		tz := MemoryStorage(config.Snappy, config.DatabaseName)
		err := tz.LoadTimezones()
		return tz, err

	} else if config.DatabaseType == "boltdb" {
		tz := BoltdbStorage(config.Snappy, config.DatabaseName, config.Encoding)
		err := tz.LoadTimezones()
		return tz, err
	}
	return &Memory{}, errors.New(errDatabaseTypeUknown)
}

func TimezonesFromGeoJSON(filename string) ([]Timezone, error) {
	start_decode := time.Now()
	fmt.Println("Building Timezone Database from: ", filename)
	var timeZones []Timezone
	file, err := os.Open(filename)
	if err != nil {
		return timeZones, err
	}
	dec := json.NewDecoder(file)

	for dec.More() {
		var js TimezoneGeoJSON

		err := dec.Decode(&js)
		if err != nil {
			return timeZones, err
		}
		for _, tz := range js.Features {
			t := Timezone{Tzid: tz.Properties.Tzid}
			switch tz.Geometry.Item {
			case "Polygon":
				t.decodePolygons(tz.Geometry.Coordinates)
			case "MultiPolygon":
				t.decodeMultiPolygons(tz.Geometry.Coordinates)
			}
			timeZones = append(timeZones, t)
		}
	}
	elapsed_decode := time.Since(start_decode)
	fmt.Println("GeoJSON decode took: ", elapsed_decode, " with ", len(timeZones), " Timezones loaded from GeoJSON")
	return timeZones, nil
}

func (t *Timezone) decodePolygons(polys []interface{}) { //1
	for _, points := range polys {
		p := t.newPolygon()
		for _, point := range points.([]interface{}) { //3
			p.updatePolygon(point.([]interface{}))
		}
		t.Polygons = append(t.Polygons, p)
	}
}

func (t *Timezone) decodeMultiPolygons(polys []interface{}) { //1
	for _, v := range polys {
		p := t.newPolygon()
		for _, points := range v.([]interface{}) { // 2
			for _, point := range points.([]interface{}) { //3
				p.updatePolygon(point.([]interface{}))
			}
		}
		t.Polygons = append(t.Polygons, p)
	}
}

func (t *Timezone) newPolygon() Polygon {
	return Polygon{
		Max: Coord{Lat: -90, Lon: -180},
		Min: Coord{Lat: 90, Lon: 180},
	}
}

func (p *Polygon) updatePolygon(xy []interface{}) {
	lon := float32(xy[0].(float64))
	lat := float32(xy[1].(float64))

	// Update max and min limits
	if p.Max.Lat < lat {
		p.Max.Lat = lat
	}
	if p.Max.Lon < lon {
		p.Max.Lon = lon
	}
	if p.Min.Lat > lat {
		p.Min.Lat = lat
	}
	if p.Min.Lon > lon {
		p.Min.Lon = lon
	}

	// add Coords to Polygon
	p.Coords = append(p.Coords, Coord{Lat: lat, Lon: lon})
}

func (p *Polygon) contains(queryPt Coord) bool {
	if len(p.Coords) < 3 {
		return false
	}
	in := rayIntersectsSegment(queryPt, p.Coords[len(p.Coords)-1], p.Coords[0])
	for i := 1; i < len(p.Coords); i++ {
		if rayIntersectsSegment(queryPt, p.Coords[i-1], p.Coords[i]) {
			in = !in
		}
	}
	return in
}

func rayIntersectsSegment(p, a, b Coord) bool {
	return (a.Lon > p.Lon) != (b.Lon > p.Lon) &&
		p.Lat < (b.Lat-a.Lat)*(p.Lon-a.Lon)/(b.Lon-a.Lon)+a.Lat
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Allocated Memory = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotal Allocated Memory = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSystem Memory = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumber of GC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
