package main
import (
	"os"
	"time"
	"log"
	"errors"
	"fmt"
	"encoding/binary"
	"runtime"
	"encoding/json"

	bolt "go.etcd.io/bbolt"
)

type TimezoneInterface interface {
	//CreatePolygonIndex() 			[]PolygonIndex
	//LoadPolygonIndex()  			[]PolygonIndex
	CreateTimezones(dbFilename string, jsonFilename string) (error) 
	LoadTimezones(filename string)						(error)
	Query(q Coord)						(string, error)
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

type PolygonIndex struct {
	Id 			uint64 		`json:"-"`
	Tzid 		string		`json:"tzid"`
	Max		    Coord 		`json:"max"`
	Min 	    Coord 		`json:"min"`
}

type Coord struct {
	X 		float64 		`json:"x"`
	Y 		float64			`json:"y"`
} 

type Store struct { 	// Database struct
	db 		*bolt.DB
	pIndex 	[]PolygonIndex
}

type Memory struct { // Memory struct
	timezones 	[]Timezone
}

var store = Store{
	pIndex: []PolygonIndex{},
}

var memory = Memory{
	timezones: []Timezone{},
}

var tz TimezoneInterface

func main() {
	//timeZones = Install("test.json") 
	//store.InsertTimezones(timeZones)
	PrintMemUsage() 
	tz = &memory
	//err := tz.CreateTimezones("timezone.db", "combined-with-oceans.json")
	//if err != nil {
	//	log.Println(err)
	//}
	tz.LoadTimezones("timezone")
	PrintMemUsage() 

	querys := []Coord{
			{Y: 5.261417, X: -3.925778,}, // Abijan Airport
			{Y: -15.678889,X: 34.973889,}, // Blantyre Airport
			{X: -53.8825,Y: 28.0325,}, // Minsk Airport
			{Y: -25.65945, X: 28.25674,}, //lat, long
			{Y: -12.65945, X: 18.25674,}, //lat, long
	    	{Y: 41.8976, X:-87.6205},
		    {Y: 47.6897, X: -122.4023},
		    {Y: 42.7235, X:-73.6931},
		    {Y: 42.5807, X:-83.0223},
		    {Y: 36.8381, X:-84.8500},
		    {Y: 40.1674, X:-85.3583},
		    {Y: 37.9643, X:-86.7453},
		    {Y: 38.6043, X:-90.2417},
		    {Y: 41.1591, X:-104.8261}, 
		    {Y: 35.1991, X:-111.6348}, 
		    {Y: 43.1432, X:-115.6750}, 
		    {Y: 47.5886, X:-122.3382}, 
		    {Y: 58.3168, X:-134.4397}, 
		    {Y: 21.4381, X:-158.0493}, 
		    {Y: 42.7000, X:-80.0000}, 
		    {Y: 51.0036, X:-114.0161}, 
		    {Y:-16.4965, X:-68.1702}, 
		    {Y:-31.9369, X:115.8453}, 
		    {Y: 42.0000, X:-87.5000}, 
	    	{Y: 41.8976, X:-87.6205},
		    {Y: 47.6897, X: -122.4023},
		    {Y: 42.7235, X:-73.6931},
		    {Y: 42.5807, X:-83.0223},
		    {Y: 36.8381, X:-84.8500},
		    {Y: 40.1674, X:-85.3583},
		    {Y: 37.9643, X:-86.7453},
		    {Y: 38.6043, X:-90.2417},
		    {Y: 41.1591, X:-104.8261}, 
		    {Y: 35.1991, X:-111.6348}, 
		    {Y: 43.1432, X:-115.6750}, 
		    {Y: 47.5886, X:-122.3382}, 
		    {Y: 58.3168, X:-134.4397}, 
		    {Y: 21.4381, X:-158.0493}, 
		    {Y: 42.7000, X:-80.0000}, 
		    {Y: 51.0036, X:-114.0161}, 
		    {Y:-16.4965, X:-68.1702}, 
		    {Y:-31.9369, X:115.8453}, 
		    {Y: 42.0000, X:-87.5000}, 
	    	{Y: 41.8976, X:-87.6205},
		    {Y: 47.6897, X: -122.4023},
		    {Y: 42.7235, X:-73.6931},
		    {Y: 42.5807, X:-83.0223},
		    {Y: 36.8381, X:-84.8500},
		    {Y: 40.1674, X:-85.3583},
		    {Y: 37.9643, X:-86.7453},
		    {Y: 38.6043, X:-90.2417},
		    {Y: 41.1591, X:-104.8261}, 
		    {Y: 35.1991, X:-111.6348}, 
		    {Y: 43.1432, X:-115.6750}, 
		    {Y: 47.5886, X:-122.3382}, 
		    {Y: 58.3168, X:-134.4397}, 
		    {Y: 21.4381, X:-158.0493}, 
		    {Y: 42.7000, X:-80.0000}, 
		    {Y: 51.0036, X:-114.0161}, 
		    {Y:-16.4965, X:-68.1702}, 
		    {Y:-31.9369, X:115.8453}, 
		    {Y: 42.0000, X:-87.5000}, 
		}
	
	for _, query := range querys {
		start := time.Now()
		res, err := tz.Query(query)
		if err != nil {
			log.Println(err)
		}
		elapsed := time.Since(start)
		fmt.Println("Query Result: ", res, " took: ", elapsed)
	}

	tz.Close()
}

func (m *Memory)Close() {
	m.timezones = []Timezone{}
	PrintMemUsage() 
}

func (s *Store)Close() {
	defer s.db.Close()
	PrintMemUsage() 
}

func (m *Memory)LoadTimezones(filename string) (error) {
	file, err := os.Open(filename + ".json")
	if err != nil {
		return err
	}
	dec := json.NewDecoder(file)
	var tzs []Timezone
	for dec.More() {
		
		
		err := dec.Decode(&tzs)
		if err != nil {
			return err
		}
	}
	m.timezones = tzs
	return nil
}

func (s *Store)LoadTimezones(filename string) (error) {
	s.OpenDB(filename + ".db")
	// Load indexes 
	return s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("Index"))
		
		b.ForEach(func(k, v []byte) error {
			var index PolygonIndex
			err := json.Unmarshal(v, &index)
			if err != nil {
				log.Println(err)
			}
			index.Id = binary.BigEndian.Uint64(k)
			s.pIndex = append(s.pIndex, index)
			return nil
		})
		return nil
	})
}

func (m *Memory)Query(q Coord) (string, error) {
	for _, tz := range m.timezones {
		for _, p := range tz.Polygons {
			if p.Min.X < q.X && p.Min.Y < q.Y && p.Max.X > q.X && p.Max.Y > q.Y {
				if p.contains(q) {
					return tz.Tzid, nil
				}
			}
		}
	}
	return "Error", errors.New("Timezone not found")
}

func (s *Store)Query(q Coord) (string, error) {
	for _, i := range s.pIndex {
		if i.Min.X < q.X && i.Min.Y < q.Y && i.Max.X > q.X && i.Max.Y > q.Y {
			p, err := s.loadPolygon(i.Id)
			if err != nil {
				return "Error", err
			} 
			if p.contains(q) {
				return i.Tzid, nil
			}
		}
	}
	return "Error", errors.New("Timezone not found")
}

func (m *Memory)writeTimezoneJSON(dbFilename string) (error) {
	w, err := os.Create(dbFilename)
    if err != nil {
    	return err
    }
    defer w.Close()
    data, err := json.Marshal(m.timezones)
    if err != nil {
    	return err
    }
    _ , err = w.Write(data)
    return err
}

func (m *Memory)CreateTimezones(dbFilename string, jsonFilename string) (error)  {
	tzs, err := TimezonesFromGeoJSON(jsonFilename)
	if err != nil {
		return err
	}
	m.timezones = tzs
	err = m.writeTimezoneJSON(dbFilename)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store)CreateTimezones(dbFilename string, jsonFilename string) (error)  {
	s.OpenDB(dbFilename)
	tzs, err := TimezonesFromGeoJSON(jsonFilename)
	if err != nil {
		return err
	}
	for _, tz := range tzs {
		s.InsertPolygons(tz)
	}
	return nil
}

func (s *Store)InsertPolygons(tz Timezone) {
	for _, polygon := range tz.Polygons {
		s.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Polygon"))
			i := tx.Bucket([]byte("Index"))

			// Get ID number autoIncrement
			id, _ := b.NextSequence()
        	intId := int(id)

        	// Create Polygon Index
        	index := PolygonIndex{
        		Tzid: tz.Tzid,
				Max: polygon.Max,
				Min: polygon.Min,
        	}

        	// UnMarshal Polygon Index
        	bufPolygon, err := json.Marshal(polygon)
		    if err != nil {
		        return err
		    }
		    // UnMarshal Polygon
        	bufIndex, err := json.Marshal(index)
		    if err != nil {
		        return err
		    }
		    // Write Polygon Index
		    err = i.Put(itob(intId), bufIndex)
		    if err != nil {
		    	return err
		    }
		    return b.Put(itob(intId), bufPolygon)
		})
	}
}

func (s *Store)loadPolygon(id uint64) (Polygon, error) {
	var polygon Polygon
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Polygon"))
		v := b.Get(itob(int(id)))
	
		return json.Unmarshal(v, &polygon)
	})
	return polygon, err
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func (s *Store)OpenDB(path string) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Println(err)
	  //return err
	}
	db.Update(func(tx *bolt.Tx) error {
	_, err := tx.CreateBucket([]byte("Index"))
	_, err = tx.CreateBucket([]byte("Polygon"))
	if err != nil {
		log.Println(err)
		//return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	s.db = db
	//defer db.Close()
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

func Install(filename string) []Timezone {
	start_decode := time.Now()
	timeZones, err := TimezonesFromGeoJSON(filename)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Timezones Loaded: ", len(timeZones))
	//log.Println(len)
	elapsed_decode := time.Since(start_decode)
	log.Printf("Timezone Decode took %s", elapsed_decode)

	
	return timeZones
}

func TimezonesFromGeoJSON(filename string) ([]Timezone, error) {
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
	return timeZones, nil
}

func (t *Timezone)decodePolygons(polys []interface{}) { //1
	for _, points := range polys {
		p := t.newPolygon()
		for _, point := range points.([]interface{}) { //3
			p.updatePolygon(point.([]interface{})) 
		}
		t.addPolygon(p)
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
		t.addPolygon(p)
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

func (t *Timezone)addPolygon(p Polygon) {
	t.Polygons = append(t.Polygons, p)
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




