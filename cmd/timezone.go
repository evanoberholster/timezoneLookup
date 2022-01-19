// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	timezone "github.com/evanoberholster/timezoneLookup"
)

var (
	// TODO: benchmark     = flag.Bool("benchmark", false, "benchmark: runs a benchmark script")
	search        = flag.Bool("search", false, "search: for the timezone for the given latitude and longitude")
	build         = flag.Bool("build", false, "build: is used to download and build timezone data")
	url           = flag.String("url", timezone.DefaultURL, "Url for data source as a zipfile default:"+timezone.DefaultURL)
	dbFilename    = flag.String("db", "timezone.data", "filename where timezone polygon data will be stored")
	cacheFilename = flag.String("cache", "/tmp/geoJSON.zip", "cache directory for downloaded zipfile")
)

func main() {
	flag.Parse()
	timezone.Verbose(true)
	if *build {
		fmt.Println("Building timezone database")
		if err := downloadAndBuild(); err != nil {
			log.Fatalln(err)
		}
	} else if *search {
		args := flag.Args()
		if len(args) != 2 {
			fmt.Println("usage: -search 'Latitude' 'Longitude'")
			fmt.Println("example: -search 10.34343 -96.3444")
			return
		} else {

			lat, err1 := strconv.ParseFloat(args[0], 64)
			lng, err2 := strconv.ParseFloat(args[1], 64)
			if err1 != nil || err2 != nil {
				fmt.Println("unable to parse: search", args[0], args[1])
				fmt.Println("usage: -search 'Latitude' 'Longitude'")
				fmt.Println("example: -search 10.34343 -96.3444")
				return
			}
			start := time.Now()
			fmt.Println("Searching timezone database")
			res, err := searchTimezone(lat, lng)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Latitude:", lat, "Longitude:", lng, "Timezone:", res.Name, "Lookup time:", res.Elapsed)
			fmt.Println("Search took:", time.Since(start))
		}

	} else {
		fmt.Println("Please choose one of the following options:")
		fmt.Println("\t", flag.Lookup("build").Usage)
		fmt.Println("\t\t", "example: timezone -build")
		fmt.Println("\t", flag.Lookup("search").Usage)
		fmt.Println("\t\t", "example: timezone -search 10.34343 -96.3444")
	}

}

func searchTimezone(lat, lng float64) (timezone.Result, error) {
	var tzc timezone.Timezonecache
	f, err := os.OpenFile(*dbFilename, syscall.O_RDWR, 0644)
	if err != nil {
		return timezone.Result{}, err
	}
	defer f.Close()
	if err = tzc.Load(f); err != nil {
		return timezone.Result{}, err
	}
	defer tzc.Close()

	return tzc.Search(lat, lng), nil

	return timezone.Result{}, err
}

func downloadAndBuild() (err error) {
	var tzc timezone.Timezonecache
	var total int
	err = timezone.ImportZipFile(*cacheFilename, *url, func(tz timezone.Timezone) error {
		total += len(tz.Polygons)
		tzc.AddTimezone(tz)
		return nil
	})
	if err != nil {
		return err
	}
	if err = tzc.Save(*dbFilename); err != nil {
		return err
	}
	fmt.Println("Polygons added:", total)
	fmt.Println("Saved Timezone data to:", *dbFilename)
	return nil
}
