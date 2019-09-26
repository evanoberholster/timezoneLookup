package main
import (
	"time"
	"fmt"
	timezone "github.com/evanoberholster/timezoneLookup"
)

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

	querys := []timezone.Coord{
			{Lat: 5.261417, Lon: -3.925778,}, // Abijan Airport
			{Lat: -15.678889,Lon: 34.973889,}, // Blantyre Airport
			{Lat: -12.65945, Lon: 18.25674,}, 
	    	{Lat: 41.8976, Lon:-87.6205},
		    {Lat: 47.6897, Lon: -122.4023},
		    {Lat: 42.7235, Lon:-73.6931},
		    {Lat: 42.5807, Lon:-83.0223},
		    {Lat: 36.8381, Lon:-84.8500},
		    {Lat: 40.1674, Lon:-85.3583},
		    {Lat: 37.9643, Lon:-86.7453},
		    {Lat: 38.6043, Lon:-90.2417},
		    {Lat: 41.1591, Lon:-104.8261}, 
		    {Lat: 35.1991, Lon:-111.6348}, 
		    {Lat: 43.1432, Lon:-115.6750}, 
		    {Lat: 47.5886, Lon:-122.3382}, 
		    {Lat: 58.3168, Lon:-134.4397}, 
		    {Lat: 21.4381, Lon:-158.0493}, 
		    {Lat: 42.7000, Lon:-80.0000}, 
		    {Lat: 51.0036, Lon:-114.0161}, 
		    {Lat:-16.4965, Lon:-68.1702}, 
		    {Lat:-31.9369, Lon:115.8453}, 
		    {Lat: 42.0000, Lon:-87.5000}, 
	    	{Lat: 41.8976, Lon:-87.6205},
		    {Lat: 47.6897, Lon: -122.4023},
		    {Lat: 42.7235, Lon:-73.6931},
		    {Lat: 42.5807, Lon:-83.0223},
		    {Lat: 36.8381, Lon:-84.8500},
		    {Lat: 40.1674, Lon:-85.3583},
		    {Lat: 37.9643, Lon:-86.7453},
		    {Lat: 38.6043, Lon:-90.2417},
		    {Lat: 41.1591, Lon:-104.8261}, 
		    {Lat: 35.1991, Lon:-111.6348}, 
		    {Lat: 43.1432, Lon:-115.6750}, 
		    {Lat: 47.5886, Lon:-122.3382}, 
		    {Lat: 58.3168, Lon:-134.4397}, 
		    {Lat: 21.4381, Lon:-158.0493}, 
		    {Lat: 42.7000, Lon:-80.0000}, 
		    {Lat: 51.0036, Lon:-114.0161}, 
		    {Lat:-16.4965, Lon:-68.1702}, 
		    {Lat:-31.9369, Lon:115.8453}, 
		    {Lat: 42.0000, Lon:-87.5000}, 
	    	{Lat: 41.8976, Lon:-87.6205},
		    {Lat: 47.6897, Lon: -122.4023},
		    {Lat: 42.7235, Lon:-73.6931},
		    {Lat: 42.5807, Lon:-83.0223},
		    {Lat: 36.8381, Lon:-84.8500},
		    {Lat: 40.1674, Lon:-85.3583},
		    {Lat: 37.9643, Lon:-86.7453},
		    {Lat: 38.6043, Lon:-90.2417},
		    {Lat: 41.1591, Lon:-104.8261}, 
		    {Lat: 35.1991, Lon:-111.6348}, 
		    {Lat: 43.1432, Lon:-115.6750}, 
		    {Lat: 47.5886, Lon:-122.3382}, 
		    {Lat: 58.3168, Lon:-134.4397}, 
		    {Lat: 21.4381, Lon:-158.0493}, 
		    {Lat: 42.7000, Lon:-80.0000}, 
		    {Lat: 51.0036, Lon:-114.0161}, 
		    {Lat:-16.4965, Lon:-68.1702}, 
		    {Lat:-31.9369, Lon:115.8453}, 
		    {Lat: 42.0000, Lon:-87.5000}, 
		}

	var times []int64
	var total int64
	
	for _, query := range querys {
		start := time.Now()
		res, err := tz.Query(query)
		if err != nil {
			fmt.Println(err)
		}
		elapsed := time.Since(start)
		fmt.Println("Query Result: ", res, " took: ", elapsed)
		times = append(times, elapsed.Nanoseconds())
		total +=  elapsed.Nanoseconds()
	}

	fmt.Println("Average time per query: ", time.Duration(total/int64(len(times))))
	tz.Close()
	timezone.PrintMemUsage() 
}
