package main

import (
	"flag"
	timezone "github.com/evanoberholster/timezoneLookup"
	"log"
)

var (
	snappy       = flag.Bool("snappy", true, "Use Snappy compression (true/false)")
	jsonFilename = flag.String("json", "combined-with-oceans.json", "GEOJSON Filename")
	dbFilename   = flag.String("db", "timezone", "Destination database filename")
	storageType  = flag.String("type", "boltdb", "Storage: boltdb or memory")
	encoding     = flag.String("encoding", "msgpack", "BoltDB encoding type: json or msgpack")
)

func main() {
	flag.Parse()

	if *dbFilename == "" || *jsonFilename == "" {
		log.Printf("Options:\n\t -snappy=true\t Use Snappy compression\n\t -json=filename\t GEOJSON filename \n\t -db=filename\t Database destination\n\t -type=boltdb\t Type of Storage (boltdb or memory) ")
	} else {
		var tz timezone.TimezoneInterface
		if *storageType == "memory" {
			tz = timezone.MemoryStorage(*snappy, *dbFilename)
		} else if *storageType == "boltdb" {
			tz = timezone.BoltdbStorage(*snappy, *dbFilename, *encoding)
		} else {
			log.Println("\"-db\" No database type specified")
			return
		}

		if *jsonFilename != "" {
			err := tz.CreateTimezones(*jsonFilename)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			log.Println("\"-json\" No GeoJSON source file specified")
			return
		}

		tz.Close()
	}

}
