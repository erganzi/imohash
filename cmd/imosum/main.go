// imosum is a sample application using imohash. It will calculate and report
// file hashes in a format similar to md5sum, etc.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/erganzi/imohash"
)

func main() {
	FilenamePtr := flag.String("filename", "foo.txt", "Filename")
	SampleSizePtr := flag.Int("samplesize", 16384, "SampleSize,Unit:bytes")
	SampleThresholdPtr := flag.Int("samplethreshold", 131072, "SampleThreshold,Unit:bytes")
	flag.Parse()

	//fmt.Println("filename:", *FilenamePtr)
	//fmt.Println("samplesize:", *SampleSizePtr)
	//fmt.Println("samplethreshold:", *SampleThresholdPtr)

	hash, err := imohash.SumFile2(*FilenamePtr, *SampleSizePtr, *SampleThresholdPtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%016x\n", hash)

}
