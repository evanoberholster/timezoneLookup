# Timezone lookup library

This is a lookup library for lat and long to Timezone.

## Data
Data should be obtained via GEOJSON, I recommend: https://github.com/evansiroky/timezone-boundary-builder

## Example
 1. Download data from: https://github.com/evansiroky/timezone-boundary-builder/releases
 2. go get github.com/evanoberholster/timezoneLookup
 3. go build github.com/evanoberholster/timezoneLookup/cmd/timezone.go
 4. timezone -j "jsonfilename" -db timezone
 5. go run example.go

```golang

package main
import (
	"time"
	"fmt"
	timezone "github.com/evanoberholster/timezoneLookup"
)

var tz timezone.TimezoneInterface

func main() {
	//tz = timezone.BoltdbStorage(timezone.WithSnappy, "timezone")
	tz = timezone.MemoryStorage(timezone.WithSnappy, "timezone")

	err := tz.LoadTimezones()
	if err != nil {
		fmt.Println(err)
	}

	querys := []timezone.Coord{
			{Y: 5.261417, X: -3.925778,}, // Abijan Airport
			{Y: -15.678889,X: 34.973889,}, // Blantyre Airport
		    {Y:-16.4965, X:-68.1702}, 
		    {Y:-31.9369, X:115.8453}, 
		    {Y: 42.0000, X:-87.5000}, 
		}
	
	for _, query := range querys {
		start := time.Now()
		res, err := tz.Query(query)
		if err != nil {
			fmt.Println(err)
		}
		elapsed := time.Since(start)
		fmt.Println("Query Result: ", res, " took: ", elapsed)
	}

	tz.Close()
}

```