package main
import (
	"fmt"
	timezone "github.com/evanoberholster/timezoneLookup"
)
var tz timezone.TimezoneInterface

func main() {
	// BoltDB storage with Snappy Compression and MsgPack encoding
	tz = timezone.BoltdbStorage(timezone.WithSnappy, "timezone", "msgpack")

	// BoltDB storage with Snappy Compression and json encoding
	//tz = timezone.BoltdbStorage(timezone.WithSnappy, "timezone", "json")

	// Storage in JSON file with Snappy Compression loaded and queried from Memory
	//tz = timezone.MemoryStorage(timezone.WithSnappy, "memory")

	err := tz.LoadTimezones()
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
