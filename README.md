# Timezone lookup
Is a timezone lookup API for GPS Coordinates to Timezone based on GeoJSON Polygon and Multipolygon files.

### Source Data
Source files can be obtained at https://github.com/evansiroky/timezone-boundary-builder

### Performance
Performance is significantly improved between release v1 and release v2 due to the underlying data structure.
Timezone database is approxiamtely 50mb in size and lookups range between 50 - 200 microseconds.
Recent verion uses an RTree as well as a Memory mapped timezone database for reduced latency and increased throughput.

[]Benchmarks still need to be written.

### Authors
Appreciate all who have contributed with Pull Requests and Issues. We eagerly welcome suggestions and PRs.

Rtree source design by Josh Backer [tidwall](https://github.com/tidwall/geoindex)

### Example

Build program
```
go build -o timezone cmd/main.go 
```

Download datasource and build timezone database ~50mb
```
./timezone -build
```

Test query for San Fransisco, United States (Etc/GMT+8)
```
./timezone -search -lat=37.7749 -lng=-122.4194
```

### Release V2.0 and forward 
Based on custom backing that loads data as memory mapped data.

```golang

package main

import (
	"fmt"

	timezone "github.com/evanoberholster/timezoneLookup/v2"
)

func main() {
	var tzc timezone.Timezonecache
	f, err := os.Open("timezone.data")
	if err != nil {
		return timezone.Result{}, err
	}
	defer f.Close()
	if err = tzc.Load(f); err != nil {
		return timezone.Result{}, err
	}
	defer tzc.Close()

	lat := 37.7749
	lng := -122.4194

	result, _ := tzc.Search(lat, lng)
	fmt.Println(result)
}

```



### Release V1.0 and prior
Based on JSON, MsgPack, Protobuf or Cap'n'Proto encodings with boltDB backend

#### Intructions for V1.0
 1. Download data from: https://github.com/evansiroky/timezone-boundary-builder/releases
 2. go get github.com/evanoberholster/timezoneLookup@v1.0.0
 3. go build github.com/evanoberholster/timezoneLookup/cmd/timezone.go
 4. timezone -json "jsonfilename" -db=timezone -type=(memory or boltdb)
 5. go run example.go

```golang
package main
import (
	"fmt"
	timezone "github.com/evanoberholster/timezoneLookup@v1.0.0"
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
			Lat: 5.261417, Lon: -3.925778,})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Query Result: ", res)

	tz.Close()
}
```
