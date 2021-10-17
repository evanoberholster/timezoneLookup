// Copyright 2018 Evan Oberholster.
//
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	timezone "github.com/evanoberholster/timezoneLookup"
)

var (
	snappy       = flag.Bool("snappy", true, "Use Snappy compression (true/false)")
	jsonFilename = flag.String("json", "combined-with-oceans.json", "GEOJSON Filename")
	dbFilename   = flag.String("db", "timezone", "Destination database filename")
	storageType  = flag.String("type", "boltdb", "Storage: boltdb or memory")
	jsonURL      = flag.String("json-url", "https://github.com/evansiroky/timezone-boundary-builder/releases/download/2020d/timezones-with-oceans.geojson.zip", "Download GeoJSON file from here if not exist")
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
	if _, err := os.Stat(*jsonFilename); err != nil && os.IsNotExist(err) {
		log.Println("Downloading " + *jsonURL)
		resp, err := http.Get(*jsonURL)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			return err
		}
		sr, err := zr.Open("combined-with-oceans.json")
		if err != nil {
			return err
		}
		fh, err := os.Create(*jsonFilename)
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(fh.Name()) }()
		defer fh.Close()
		if _, err = io.Copy(fh, sr); err != nil {
			return err
		}
		if err := fh.Close(); err != nil {
			return err
		}
	}
	err := tz.CreateTimezones(*jsonFilename)
	if err != nil {
		return err
	}
	tz.Close()
	return nil
}
