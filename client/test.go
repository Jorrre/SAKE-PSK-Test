package main

import (
	"SAKE-TLS_Test/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run test.go <serverAddr> <runTime in seconds>. Flags: -r enables session resumption")
		return
	}
	rt, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid run time %s", os.Args[2])
		return
	}
	if rt < 0 {
		fmt.Printf("Negative run time %d", rt)
	}
	runTime := time.Duration(rt) * time.Second

	if len(os.Args) == 4 && os.Args[3] == "-r" {
		runTests(os.Args[1], true, runTime)
	} else {
		runTests(os.Args[1], false, runTime)
	}
}

func runTests(serverAddr string, useRes bool, runTime time.Duration) {
	log.SetFlags(log.Lmicroseconds)
	parallelClients := []int{1, 2, 4, 6, 8, 10}

	var results []float64
	for _, clients := range parallelClients {
		result := runTest(clients, useRes, serverAddr, runTime)
		results = append(results, result)
	}
	log.Printf("**************** Results, %d seconds ****************", int(runTime.Seconds()))
	utils.LogResult(results)
}

func runTest(clients int, useRes bool, serverAddr string, runTime time.Duration) float64 {
	if useRes {
		log.Printf("================ Running TLS-PSK test with %d client(s) ================", clients)
	} else {
		log.Printf("================ Running full TLS handshake test with %d client(s) ================", clients)
	}

	resultChan := make(chan int, clients)
	var wg sync.WaitGroup

	log.Printf("Starting %d client(s)...\n", clients)
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go Client(serverAddr, useRes, runTime, resultChan, &wg)
	}
	log.Println("All clients up and running")
	log.Printf("Performing handshakes for %d seconds...", int(runTime.Seconds()))

	go func() {
		wg.Wait()
		close(resultChan)
		log.Println("-------- Test complete --------")
	}()

	total := 0
	for sum := range resultChan {
		total += sum
	}
	log.Printf("Total number of handshakes: %d", total)
	handshakesPerSec := float64(total) / runTime.Seconds()
	log.Printf("Handshakes per second: %s", strconv.FormatFloat(handshakesPerSec, 'f', -1, 64))

	println()
	return handshakesPerSec
}
