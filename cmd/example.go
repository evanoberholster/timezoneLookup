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
			Lat: 5.261417, Lon: -3.925778,})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Query Result: ", res)

	tz.Close()
}
