package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: go run . <filepath> <concurrency>")
	}

	filepath := os.Args[1]

	concurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("cannot parse the concurrency number: %s", err)
	}

	// Allocate a bitmap large enough to handle all possible IPv4 addresses (2 ^ 32)
	bitmap := NewBitmap(ipV4AddressesCount)

	ipDataset := NewIPDataset(filepath, concurrency)

	timeStart := time.Now()

	uniqueIPCount, err := ipDataset.CountUniqueIPs(bitmap)
	if err != nil {
		log.Fatalf("cannot count unique IPs: %s", err)
	}

	log.Printf("count of unique IP addresses: %d\n", uniqueIPCount)
	log.Printf("execution time: %s\n", time.Since(timeStart))
}
