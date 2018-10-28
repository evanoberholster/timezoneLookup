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
			{Lon: 5.261417, Lat: -3.925778,}, // Abijan Airport
			{Lon: -15.678889,Lat: 34.973889,}, // Blantyre Airport
			{Lon: -12.65945, Lat: 18.25674,}, 
	    	{Lon: 41.8976, Lat:-87.6205},
		    {Lon: 47.6897, Lat: -122.4023},
		    {Lon: 42.7235, Lat:-73.6931},
		    {Lon: 42.5807, Lat:-83.0223},
		    {Lon: 36.8381, Lat:-84.8500},
		    {Lon: 40.1674, Lat:-85.3583},
		    {Lon: 37.9643, Lat:-86.7453},
		    {Lon: 38.6043, Lat:-90.2417},
		    {Lon: 41.1591, Lat:-104.8261}, 
		    {Lon: 35.1991, Lat:-111.6348}, 
		    {Lon: 43.1432, Lat:-115.6750}, 
		    {Lon: 47.5886, Lat:-122.3382}, 
		    {Lon: 58.3168, Lat:-134.4397}, 
		    {Lon: 21.4381, Lat:-158.0493}, 
		    {Lon: 42.7000, Lat:-80.0000}, 
		    {Lon: 51.0036, Lat:-114.0161}, 
		    {Lon:-16.4965, Lat:-68.1702}, 
		    {Lon:-31.9369, Lat:115.8453}, 
		    {Lon: 42.0000, Lat:-87.5000}, 
	    	{Lon: 41.8976, Lat:-87.6205},
		    {Lon: 47.6897, Lat: -122.4023},
		    {Lon: 42.7235, Lat:-73.6931},
		    {Lon: 42.5807, Lat:-83.0223},
		    {Lon: 36.8381, Lat:-84.8500},
		    {Lon: 40.1674, Lat:-85.3583},
		    {Lon: 37.9643, Lat:-86.7453},
		    {Lon: 38.6043, Lat:-90.2417},
		    {Lon: 41.1591, Lat:-104.8261}, 
		    {Lon: 35.1991, Lat:-111.6348}, 
		    {Lon: 43.1432, Lat:-115.6750}, 
		    {Lon: 47.5886, Lat:-122.3382}, 
		    {Lon: 58.3168, Lat:-134.4397}, 
		    {Lon: 21.4381, Lat:-158.0493}, 
		    {Lon: 42.7000, Lat:-80.0000}, 
		    {Lon: 51.0036, Lat:-114.0161}, 
		    {Lon:-16.4965, Lat:-68.1702}, 
		    {Lon:-31.9369, Lat:115.8453}, 
		    {Lon: 42.0000, Lat:-87.5000}, 
	    	{Lon: 41.8976, Lat:-87.6205},
		    {Lon: 47.6897, Lat: -122.4023},
		    {Lon: 42.7235, Lat:-73.6931},
		    {Lon: 42.5807, Lat:-83.0223},
		    {Lon: 36.8381, Lat:-84.8500},
		    {Lon: 40.1674, Lat:-85.3583},
		    {Lon: 37.9643, Lat:-86.7453},
		    {Lon: 38.6043, Lat:-90.2417},
		    {Lon: 41.1591, Lat:-104.8261}, 
		    {Lon: 35.1991, Lat:-111.6348}, 
		    {Lon: 43.1432, Lat:-115.6750}, 
		    {Lon: 47.5886, Lat:-122.3382}, 
		    {Lon: 58.3168, Lat:-134.4397}, 
		    {Lon: 21.4381, Lat:-158.0493}, 
		    {Lon: 42.7000, Lat:-80.0000}, 
		    {Lon: 51.0036, Lat:-114.0161}, 
		    {Lon:-16.4965, Lat:-68.1702}, 
		    {Lon:-31.9369, Lat:115.8453}, 
		    {Lon: 42.0000, Lat:-87.5000}, 
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
