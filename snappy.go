package main

import (
	"flag"
	"io"
	"os"
	"log"
	"fmt"
	"time"

	"github.com/golang/snappy"
)

var (
	decode = flag.Bool("d", false, "decode")
	encode = flag.Bool("e", false, "encode")
)


func main() {
	file, err := os.Open("combined-with-oceans.json")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fi, err := os.Create("test.snap")
    if err != nil {
    	log.Println(err)
    }
    defer fi.Close()
    start_encode := time.Now()
	snap := snappy.NewBufferedWriter(fi)
	_, err = io.Copy(snap, file)
	  if err != nil {
	    log.Fatal(err)
	  }
	defer snap.Close()
	elapsed_encode := time.Since(start_encode)
	fmt.Println("Snappy Endcode took: ", elapsed_encode)
}