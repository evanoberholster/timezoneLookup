// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"flag"
	"log"

	timezone "github.com/evanoberholster/timezoneLookup"
)

var (
	snappy       = flag.Bool("snappy", true, "Use Snappy compression (true/false)")
	jsonFilename = flag.String("json", "combined-with-oceans.json", "GEOJSON Filename")
	dbFilename   = flag.String("db", "timezone", "Destination database filename")
	storageType  = flag.String("type", "boltdb", "Storage: boltdb or memory")
	encoding     = flag.String("encoding", "msgpack", "BoltDB encoding type: json or msgpack")
)

func main() {
	if err := Main(); err != nil {
		log.Fatalln(err)
	}
}
func Main() error {
	flag.Parse()

	if *dbFilename == "" || *jsonFilename == "" {
		log.Printf("Options:\n\t -snappy=true\t Use Snappy compression\n\t -json=filename\t GEOJSON filename \n\t -db=filename\t Database destination\n\t -type=boltdb\t Type of Storage (boltdb or memory) ")
		return nil
	}
	var tz timezone.TimezoneInterface
	if *storageType == "memory" {
		tz = timezone.MemoryStorage(*snappy, *dbFilename)
	} else if *storageType == "boltdb" {
		tz = timezone.BoltdbStorage(*snappy, *dbFilename, *encoding)
	} else {
		return errors.New("\"-db\" No database type specified")
	}

	if *jsonFilename == "" {
		return errors.New("\"-json\" No GeoJSON source file specified")
	}
	err := tz.CreateTimezones(*jsonFilename)
	if err != nil {
		return err
	}
	tz.Close()
	return nil
}
