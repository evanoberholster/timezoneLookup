# Timezone lookup library
This is a lookup API for GPS Coordinates to Timezone based on GeoJSON files. 

- Support added for MsgPack encoding for faster query.

## Data
Data should be obtained via GeoJSON, I recommend: https://github.com/evansiroky/timezone-boundary-builder

## Performance
```
 go run benchmark.go
```

### Performance BoltDB
With Snappy Compression and data loaded from Boltdb

(Msgpack) Average time per query:  ~18.015757ms 

(Msgpack) Database file size: 89M

(Json) Average time per query:  ~38.356527ms

(Json) Database file size: 70M

```
Alloc = 4 MiB	TotalAlloc = 232 MiB	Sys = 68 MiB	NumGC = 88
```

### Performance Memory
With Snappy Compression and data loaded from memory

Average time per query:  ~87.982µs

Database file size: 65M 

```
Alloc = 550 MiB	TotalAlloc = 995 MiB	Sys = 601 MiB	NumGC = 10
```

## Example
 1. Download data from: https://github.com/evansiroky/timezone-boundary-builder/releases
 2. go get github.com/evanoberholster/timezoneLookup
 3. go build github.com/evanoberholster/timezoneLookup/cmd/timezone.go
 4. timezone -json "jsonfilename" -db=timezone -type=(memory or boltdb)
 5. go run example.go


```golang
package main
import (
	"fmt"
	timezone "github.com/evanoberholster/timezoneLookup"
)
var tz timezone.TimezoneInterface

func main() {
	tz, err := timezone.LoadTimezones(timezone.Config{
											DatabaseType:"boltdb", // memory or boltdb
											DatabaseName:"timezone", // Name without suffix
											Snappy: true,
											Encoding: "msgpack", // json or msgpack
										})
	if err != nil {
		fmt.Println(err)
	}

	res, err := tz.Query(timezone.Coord{
			Lon: 5.261417, Lat: -3.925778,})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Query Result: ", res)

	tz.Close()
}
```