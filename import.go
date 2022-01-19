package timezoneLookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/evanoberholster/timezoneLookup/geo"
	"github.com/klauspost/compress/zip"
)

var (
	verbose bool = false
)

// Verbose sets the import function to verbose
func Verbose(v bool) {
	verbose = v
}

const (
	DefaultURL = "https://github.com/evansiroky/timezone-boundary-builder/releases/download/2020d/timezones-with-oceans.geojson.zip"
)

// ImportZipFile imports a url and saves it with the following filename. The iter function is run on the zip file.
func ImportZipFile(cache string, url string, iter func(tz Timezone) error) (err error) {
	start := time.Now()
	if _, err := os.Stat(cache); errors.Is(err, os.ErrNotExist) {
		if verbose {
			fmt.Println("Caching url:", url, "to:", cache)
		}
		if err = fetchAndCacheFile(cache, url); err != nil {
			return err
		}
		if verbose {
			fmt.Println("Time to download timezone JSON:", time.Since(start))
			start = time.Now()
		}
	}
	if verbose {
		fmt.Println("Loading cache:", cache)
	}
	if !strings.EqualFold(cache[len(cache)-4:], ".zip") {
		return errors.New("error not a zip file")
	}
	var zr *zip.ReadCloser
	if zr, err = zip.OpenReader(cache); err != nil {
		return err
	}
	defer zr.Close()
	for _, v := range zr.File {
		if strings.EqualFold(".json", v.Name[len(v.Name)-5:]) {
			decodeJSON(v, iter)
		}
	}
	if verbose {
		fmt.Println("Time to process timezones:", time.Since(start))
	}

	return nil
}

func fetchAndCacheFile(filename string, url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	if n != resp.ContentLength {
		fmt.Println(n, resp.ContentLength)
	}

	return
}

func decodeJSON(f *zip.File, iter func(tz Timezone) error) (err error) {
	var rc io.ReadCloser
	if rc, err = f.Open(); err != nil {
		return err
	}
	defer rc.Close()

	dec := json.NewDecoder(rc)

	var token json.Token
	for dec.More() {
		if token, err = dec.Token(); err != nil {
			break
		}
		if t, ok := token.(string); ok && t == "features" {
			if token, err = dec.Token(); err == nil && token.(json.Delim) == '[' {
				return decodeFeatures(dec, iter) // decode features
			}
		}
	}
	return errors.New("error no features found")
}

func decodeFeatures(dec *json.Decoder, fn func(tz Timezone) error) error {
	var f GeoJSONFeature
	var err error

	for dec.More() {
		if err = dec.Decode(&f); err != nil {
			return err
		}
		var pp []geo.Polygon
		switch f.Geometry.Item {
		case "Polygon":
			pp = decodePolygons(f.Geometry.Coordinates)
		case "MultiPolygon":
			pp = decodeMultiPolygons(f.Geometry.Coordinates)
		}
		if err = fn(Timezone{Name: f.Properties.Tzid, Polygons: pp}); err != nil {
			return err
		}
	}

	return nil
}

// decodePolygons
// GeoJSON Spec https://geojson.org/geojson-spec.html
// Coordinates: [Longitude, Latitude]
func decodePolygons(polygons []interface{}) []geo.Polygon {
	var pp []geo.Polygon
	for _, points := range polygons {
		p := geo.NewPolygon()
		for _, i := range points.([]interface{}) {
			if latlng, ok := i.([]interface{}); ok {
				p.AddVertex(geo.NewLatLng(latlng[1].(float64), latlng[0].(float64)))
			}
		}
		pp = append(pp, p)
	}
	return pp
}

// decodeMultiPolygons
// GeoJSON Spec https://geojson.org/geojson-spec.html
// Coordinates: [Longitude, Latitude]
func decodeMultiPolygons(polygons []interface{}) []geo.Polygon {
	var pp []geo.Polygon
	for _, v := range polygons {
		p := geo.NewPolygon()
		for _, points := range v.([]interface{}) { // 2
			for _, i := range points.([]interface{}) {
				if latlng, ok := i.([]interface{}); ok {
					p.AddVertex(geo.NewLatLng(latlng[1].(float64), latlng[0].(float64)))
				}
			}
		}
		pp = append(pp, p)
	}
	return pp
}

// Timezone
type Timezone struct {
	Name     string
	Polygons []geo.Polygon `json:"polygons"`
}

// GeoJSONFeature
type GeoJSONFeature struct {
	Type       string `json:"type"`
	Properties struct {
		Tzid string `json:"tzid"`
	} `json:"properties"`
	Geometry struct {
		Item        string        `json:"type"`
		Coordinates []interface{} `json:"coordinates"`
	} `json:"geometry"`
}
