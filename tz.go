package main
import (
	"os"
	"time"
	"log"
	"fmt"
	"runtime"
	"encoding/json"
)

	
type TimeZoneGeoJSON struct {
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

type TimeZone struct {
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

func main() {
	timeZones := Install("test.json") 

	querys := []Coord{
			{Y: 5.261417, X: -3.925778,}, // Abijan Airport
			{Y: -15.678889,X: 34.973889,}, // Blantyre Airport
			{X: -53.8825,Y: 28.0325,}, // Minsk Airport
			{Y: -25.65945, X: 28.25674,}, //lat, long
			{Y: -1.65945, X: 18.25674,}, //lat, long
		}
	for _, query := range querys {
		start := time.Now()
		res := QueryTimeZoneWithLimits(timeZones, query)
		elapsed := time.Since(start)
		fmt.Println("Query Result: ", res, " took: ", elapsed)

		start2 := time.Now()
		res2 := QueryTimeZone(timeZones, query)
		elapsed2 := time.Since(start2)
		fmt.Println("Query Result: ", res2, " took: ", elapsed2)
	}

}

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

func Install(filename string) []TimeZone {
	start_decode := time.Now()
	timeZones, err := TransformTimeZoneJSON(filename)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Timezones Loaded: ", len(timeZones))
	//log.Println(len)
	elapsed_decode := time.Since(start_decode)
	log.Printf("TimeZone Decode took %s", elapsed_decode)

	
	return timeZones
}

func TransformTimeZoneJSON(filename string) ([]TimeZone, error) {
	var timeZones []TimeZone
	file, err := os.Open(filename)
	if err != nil {
		return timeZones, err
	}
	dec := json.NewDecoder(file)

	for dec.More() {
		var js TimeZoneGeoJSON
		
		err := dec.Decode(&js)
		if err != nil {
			return timeZones, err
		}
		for _, tz := range js.Features {
			t := TimeZone{Tzid: tz.Properties.Tzid}
			switch tz.Geometry.Item {
				case "Polygon":
					t.DecodePolygons(tz.Geometry.Coordinates)
				case "MultiPolygon":
					t.DecodeMultiPolygons(tz.Geometry.Coordinates)
			}
			timeZones = append(timeZones, t)
		}
	}
	return timeZones, nil
}

func (t *TimeZone)DecodePolygons(polys []interface{}) { //1
	p := t.newPolygon()
	for _, points := range polys {
		for _, point := range points.([]interface{}) { //3
			p.updatePolygon(point.([]interface{})) 
		}
		t.addPolygon(p)
	}
}

func (t *TimeZone)DecodeMultiPolygons(polys []interface{}) { //1
	for _, v := range polys {
		p := t.newPolygon()
		for _, points := range v.([]interface{}) { // 2
			for _, point := range points.([]interface{}) { //3
				p.updatePolygon(point.([]interface{})) 
			}
		}
		t.addPolygon(p)
	}
}

func QueryTimeZone(tzs []TimeZone, q Coord) string {
	for _, tz := range tzs {
		for _, p := range tz.Polygons {
			if p.Contains(q) {
				return tz.Tzid
			}
		}
	}
	return "Not Found"
}

func QueryTimeZoneWithLimits(tzs []TimeZone, q Coord) string {
	for _, tz := range tzs {
		for _, p := range tz.Polygons {
			if p.Min.X < q.X && p.Min.Y < q.Y && p.Max.X > q.X && p.Max.Y > q.Y {
				if p.Contains(q) {
					return tz.Tzid
				}
			}
		}
	}
	return "Not Found"
}

func (t *TimeZone)newPolygon() (Polygon) {
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

func (t *TimeZone)addPolygon(p Polygon) {
	t.Polygons = append(t.Polygons, p)
}

func (p *Polygon)Contains(queryPt Coord) bool {
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




