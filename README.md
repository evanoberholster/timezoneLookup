# Timezone lookup library

This is a lookup library for lat and long to Timezone.

## Data
Data should be obtained via GEOJSON, I recommend: https://github.com/evansiroky/timezone-boundary-builder

## Performance BoltDB
With Snappy Compression and data loaded from Boltdb
Database file size: 68M
```
Alloc = 4 MiB	TotalAlloc = 232 MiB	Sys = 68 MiB	NumGC = 88
Query Result:  Africa/Abidjan  took:  26.648672ms
Query Result:  Africa/Blantyre  took:  22.681323ms
Query Result:  Africa/Luanda  took:  38.252324ms
Query Result:  America/Chicago  took:  54.988383ms
Query Result:  America/Los_Angeles  took:  40.447983ms
Query Result:  America/New_York  took:  46.139995ms
Query Result:  America/Detroit  took:  6.125602ms
Query Result:  America/Kentucky/Monticello  took:  58.695582ms
Query Result:  America/Indiana/Indianapolis  took:  76.587206ms
Query Result:  America/Indiana/Tell_City  took:  59.79741ms
Query Result:  America/Chicago  took:  78.330179ms
Query Result:  America/Denver  took:  113.776201ms
Query Result:  America/Phoenix  took:  52.412518ms
Query Result:  America/Boise  took:  21.605222ms
Query Result:  America/Los_Angeles  took:  38.692886ms
Query Result:  America/Juneau  took:  1.677749ms
Query Result:  Pacific/Honolulu  took:  332.136µs
Query Result:  America/Toronto  took:  84.911477ms
Query Result:  America/Edmonton  took:  11.981108ms
Query Result:  America/La_Paz  took:  38.192319ms
Query Result:  Australia/Perth  took:  2.557322ms
Query Result:  America/Chicago  took:  60.662676ms
Query Result:  America/Chicago  took:  61.883292ms
Query Result:  America/Los_Angeles  took:  32.844784ms
Query Result:  America/New_York  took:  39.29996ms
Query Result:  America/Detroit  took:  5.280706ms
Query Result:  America/Kentucky/Monticello  took:  53.435807ms
Query Result:  America/Indiana/Indianapolis  took:  54.058419ms
Query Result:  America/Indiana/Tell_City  took:  60.79631ms
Query Result:  America/Chicago  took:  55.65176ms
Query Result:  America/Denver  took:  96.298662ms
Query Result:  America/Phoenix  took:  44.399702ms
Query Result:  America/Boise  took:  16.578683ms
Query Result:  America/Los_Angeles  took:  35.365181ms
Query Result:  America/Juneau  took:  943.142µs
Query Result:  Pacific/Honolulu  took:  224.829µs
Query Result:  America/Toronto  took:  67.522004ms
Query Result:  America/Edmonton  took:  11.54448ms
Query Result:  America/La_Paz  took:  32.70719ms
Query Result:  Australia/Perth  took:  2.448756ms
Query Result:  America/Chicago  took:  58.987711ms
Query Result:  America/Chicago  took:  61.025159ms
Query Result:  America/Los_Angeles  took:  31.941886ms
Query Result:  America/New_York  took:  40.119421ms
Query Result:  America/Detroit  took:  4.470533ms
Query Result:  America/Kentucky/Monticello  took:  56.199474ms
Query Result:  America/Indiana/Indianapolis  took:  58.214663ms
Query Result:  America/Indiana/Tell_City  took:  55.784476ms
Query Result:  America/Chicago  took:  53.909408ms
Query Result:  America/Denver  took:  100.816517ms
Query Result:  America/Phoenix  took:  42.115515ms
Query Result:  America/Boise  took:  16.02315ms
Query Result:  America/Los_Angeles  took:  30.808792ms
Query Result:  America/Juneau  took:  1.290399ms
Query Result:  Pacific/Honolulu  took:  225.36µs
Query Result:  America/Toronto  took:  68.879193ms
Query Result:  America/Edmonton  took:  11.929186ms
Query Result:  America/La_Paz  took:  33.506021ms
Query Result:  Australia/Perth  took:  2.657613ms
Query Result:  America/Chicago  took:  58.477893ms
```

## Performance Memory
With Snappy Compression and data loaded from memory
Database file size: 65M 
```
Alloc = 550 MiB	TotalAlloc = 995 MiB	Sys = 601 MiB	NumGC = 10
Query Result:  Africa/Abidjan  took:  48.472µs
Query Result:  Africa/Blantyre  took:  65.852µs
Query Result:  Etc/GMT+4  took:  59.427µs
Query Result:  Africa/Luanda  took:  112.587µs
Query Result:  America/Chicago  took:  143.611µs
Query Result:  America/Los_Angeles  took:  87.94µs
Query Result:  America/New_York  took:  87.868µs
Query Result:  America/Detroit  took:  18.3µs
Query Result:  America/Kentucky/Monticello  took:  134.541µs
Query Result:  America/Indiana/Indianapolis  took:  126.347µs
Query Result:  America/Indiana/Tell_City  took:  143.905µs
Query Result:  America/Chicago  took:  108.424µs
Query Result:  America/Denver  took:  195.15µs
Query Result:  America/Phoenix  took:  127.91µs
Query Result:  America/Boise  took:  41.205µs
Query Result:  America/Los_Angeles  took:  74.436µs
Query Result:  America/Juneau  took:  7.618µs
Query Result:  Pacific/Honolulu  took:  15.506µs
Query Result:  America/Toronto  took:  173.747µs
Query Result:  America/Edmonton  took:  40.256µs
Query Result:  America/La_Paz  took:  80.842µs
Query Result:  Australia/Perth  took:  18.011µs
Query Result:  America/Chicago  took:  130.723µs
Query Result:  America/Chicago  took:  152.806µs
Query Result:  America/Los_Angeles  took:  86.001µs
Query Result:  America/New_York  took:  110.615µs
Query Result:  America/Detroit  took:  16.486µs
Query Result:  America/Kentucky/Monticello  took:  138.022µs
Query Result:  America/Indiana/Indianapolis  took:  114.373µs
Query Result:  America/Indiana/Tell_City  took:  99.879µs
Query Result:  America/Chicago  took:  96.236µs
Query Result:  America/Denver  took:  199.58µs
Query Result:  America/Phoenix  took:  84.504µs
Query Result:  America/Boise  took:  34.801µs
Query Result:  America/Los_Angeles  took:  60.395µs
Query Result:  America/Juneau  took:  4.387µs
Query Result:  Pacific/Honolulu  took:  12.988µs
Query Result:  America/Toronto  took:  137.44µs
Query Result:  America/Edmonton  took:  25.974µs
Query Result:  America/La_Paz  took:  65.172µs
Query Result:  Australia/Perth  took:  13.296µs
Query Result:  America/Chicago  took:  105.497µs
Query Result:  America/Chicago  took:  106.998µs
Query Result:  America/Los_Angeles  took:  75.97µs
Query Result:  America/New_York  took:  84.493µs
Query Result:  America/Detroit  took:  15.766µs
Query Result:  America/Kentucky/Monticello  took:  132.73µs
Query Result:  America/Indiana/Indianapolis  took:  109.399µs
Query Result:  America/Indiana/Tell_City  took:  133.941µs
Query Result:  America/Chicago  took:  122.689µs
Query Result:  America/Denver  took:  196.18µs
Query Result:  America/Phoenix  took:  93.212µs
Query Result:  America/Boise  took:  37.742µs
Query Result:  America/Los_Angeles  took:  89.8µs
Query Result:  America/Juneau  took:  6.25µs
Query Result:  Pacific/Honolulu  took:  16.606µs
Query Result:  America/Toronto  took:  161.733µs
Query Result:  America/Edmonton  took:  36.757µs
Query Result:  America/La_Paz  took:  75.31µs
Query Result:  Australia/Perth  took:  15.93µs
Query Result:  America/Chicago  took:  120.317µs
```


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