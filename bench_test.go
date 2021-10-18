// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package timezoneLookup_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	timezone "github.com/evanoberholster/timezoneLookup"
)

func BenchmarkLookup(b *testing.B) {
	_ = os.MkdirAll("testdata", 0755)
	tzgo := filepath.Join("..", "cmd", "timezone.go")
	for _, e := range []string{"msgpack.snap", "msgpack", "protobuf.snap", "protobuf", "capnp.snap", "capnp", "snap.json"} {
		cfg := timezone.Config{
			DatabaseName: filepath.Join("testdata", "timezone"),
		}
		if e == "json" || e == "snap.json" {
			cfg.Snappy = strings.HasPrefix(e, "snap.")
			if _, err := os.Stat(cfg.DatabaseName + "." + e); err != nil && os.IsNotExist(err) {
				cmd := exec.Command("go", "run", tzgo, "-type=memory", "-snappy="+strconv.FormatBool(cfg.Snappy))
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
				cmd.Dir = "testdata"
				_ = cmd.Run()
			}
			cfg.DatabaseType = "memory"
		} else {
			cfg.Snappy = strings.HasSuffix(e, ".snap")
			_, err := os.Stat(cfg.DatabaseName + "." + e + ".db")
			e := strings.TrimSuffix(e, ".snap")
			if err != nil && os.IsNotExist(err) {
				cmd := exec.Command("go", "run", tzgo, "-encoding="+e, "-snappy="+strconv.FormatBool(cfg.Snappy))
				cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
				cmd.Dir = "testdata"
				_ = cmd.Run()
			}
			if cfg.Encoding, err = timezone.EncodingFromString(e); err != nil {
				b.Fatal(err)
			}
			cfg.DatabaseType = "boltdb"
		}
		b.Run(e, func(b *testing.B) {
			tz, err := timezone.LoadTimezones(cfg)
			if err != nil {
				b.Fatalf("%q: %#v: %+v", e, cfg, err)
			}
			defer tz.Close()

			benchLookup(b, tz)
		})
	}
}

func benchLookup(b *testing.B, tz timezone.TimezoneInterface) {
	querys := []timezone.Coord{
		{Lat: 5.261417, Lon: -3.925778},   // Abijan Airport
		{Lat: -15.678889, Lon: 34.973889}, // Blantyre Airport
		{Lat: -12.65945, Lon: 18.25674},
		{Lat: 41.8976, Lon: -87.6205},
		{Lat: 47.6897, Lon: -122.4023},
		{Lat: 42.7235, Lon: -73.6931},
		{Lat: 42.5807, Lon: -83.0223},
		{Lat: 36.8381, Lon: -84.8500},
		{Lat: 40.1674, Lon: -85.3583},
		{Lat: 37.9643, Lon: -86.7453},
		{Lat: 38.6043, Lon: -90.2417},
		{Lat: 41.1591, Lon: -104.8261},
		{Lat: 35.1991, Lon: -111.6348},
		{Lat: 43.1432, Lon: -115.6750},
		{Lat: 47.5886, Lon: -122.3382},
		{Lat: 58.3168, Lon: -134.4397},
		{Lat: 21.4381, Lon: -158.0493},
		{Lat: 42.7000, Lon: -80.0000},
		{Lat: 51.0036, Lon: -114.0161},
		{Lat: -16.4965, Lon: -68.1702},
		{Lat: -31.9369, Lon: 115.8453},
		{Lat: 42.0000, Lon: -87.5000},
		{Lat: 41.8976, Lon: -87.6205},
		{Lat: 47.6897, Lon: -122.4023},
		{Lat: 42.7235, Lon: -73.6931},
		{Lat: 42.5807, Lon: -83.0223},
		{Lat: 36.8381, Lon: -84.8500},
		{Lat: 40.1674, Lon: -85.3583},
		{Lat: 37.9643, Lon: -86.7453},
		{Lat: 38.6043, Lon: -90.2417},
		{Lat: 41.1591, Lon: -104.8261},
		{Lat: 35.1991, Lon: -111.6348},
		{Lat: 43.1432, Lon: -115.6750},
		{Lat: 47.5886, Lon: -122.3382},
		{Lat: 58.3168, Lon: -134.4397},
		{Lat: 21.4381, Lon: -158.0493},
		{Lat: 42.7000, Lon: -80.0000},
		{Lat: 51.0036, Lon: -114.0161},
		{Lat: -16.4965, Lon: -68.1702},
		{Lat: -31.9369, Lon: 115.8453},
		{Lat: 42.0000, Lon: -87.5000},
		{Lat: 41.8976, Lon: -87.6205},
		{Lat: 47.6897, Lon: -122.4023},
		{Lat: 42.7235, Lon: -73.6931},
		{Lat: 42.5807, Lon: -83.0223},
		{Lat: 36.8381, Lon: -84.8500},
		{Lat: 40.1674, Lon: -85.3583},
		{Lat: 37.9643, Lon: -86.7453},
		{Lat: 38.6043, Lon: -90.2417},
		{Lat: 41.1591, Lon: -104.8261},
		{Lat: 35.1991, Lon: -111.6348},
		{Lat: 43.1432, Lon: -115.6750},
		{Lat: 47.5886, Lon: -122.3382},
		{Lat: 58.3168, Lon: -134.4397},
		{Lat: 21.4381, Lon: -158.0493},
		{Lat: 42.7000, Lon: -80.0000},
		{Lat: 51.0036, Lon: -114.0161},
		{Lat: -16.4965, Lon: -68.1702},
		{Lat: -31.9369, Lon: 115.8453},
		{Lat: 42.0000, Lon: -87.5000},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := querys[i%len(querys)]
		_, err := tz.Query(query)
		if err != nil {
			fmt.Println(err)
		}
	}
	b.StopTimer()
}
