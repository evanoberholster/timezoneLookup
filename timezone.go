package timezoneLookup
import (
	"os"
	"time"
	"log"
	"fmt"
	"runtime"
	"encoding/json"
)

const (
	WithSnappy = true
	NoSnappy = false
)

type TimezoneInterface interface {
	CreateTimezones(jsonFilename string) 	(error) 
	LoadTimezones()							(error)
	Query(q Coord)							(string, error)
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
	Tzid 		string		`json:"tzid"`
	Polygons 	[]Polygon	`json:"polygons"`	
}

type Polygon struct {
	Max		    Coord 		`json:"max"`
	Min 	    Coord 		`json:"min"`
	Coords 		[]Coord 	`json:"coords"`
}

type Coord struct {
	X 		float64 		`json:"x"`
	Y 		float64			`json:"y"`
} 

var Tz TimezoneInterface

func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
func TimezonesFromGeoJSON(filename string) ([]Timezone, error) {
	start_decode := time.Now()
	fmt.Println("Loading Timezones from: ", filename)
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
	fmt.Println("Timezones Loaded: ", len(timeZones))
	elapsed_decode := time.Since(start_decode)
	log.Println("Timezones decode took: ", elapsed_decode)
	return timeZones, nil
}

func (t *Timezone)decodePolygons(polys []interface{}) { //1
	for _, points := range polys {
		p := t.newPolygon()
		for _, point := range points.([]interface{}) { //3
			p.updatePolygon(point.([]interface{})) 
		}
		t.Polygons = append(t.Polygons, p)
	}
}

func (t *Timezone)decodeMultiPolygons(polys []interface{}) { //1
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

func (t *Timezone)newPolygon() (Polygon) {
	return Polygon{
			Max: Coord{ X: -180, Y: -180, },
			Min: Coord{ X: 180, Y: 180, },
		}
}

func (p *Polygon)updatePolygon(xy []interface{}) {
	x := xy[0].(float64)
	y := xy[1].(float64)

	// Update max and min limits
	if p.Max.X < x { p.Max.X = x }
	if p.Max.Y < y { p.Max.Y = y }
	if p.Min.X > x { p.Min.X = x }
	if p.Min.Y > y { p.Min.Y = y }

	// add Coords to Polygon
	p.Coords = append(p.Coords, Coord{X:x, Y:y})
}

func (p *Polygon)contains(queryPt Coord) bool {
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
    return (a.Y > p.Y) != (b.Y > p.Y) &&
        p.X < (b.X-a.X)*(p.Y-a.Y)/(b.Y-a.Y)+a.X
}




