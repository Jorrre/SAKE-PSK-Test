package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

const maxClients = 10
const runTime = time.Second * 60
const serverIP = "127.0.0.1"

var serverPort = 2208

var serverAddress = ""

func main() {
	log.SetFlags(log.Lmicroseconds)
	for clients := 1; clients <= maxClients; clients++ {
		for _, useRes := range []bool{true, false} {
			runTest(clients, useRes)
		}
	}
}

func runTest(clients int, useRes bool) int {
	if useRes {
		log.Printf("================ Running TLS-PSK test with %d client(s) ================", clients)
	} else {
		log.Printf("================ Running full handshake test with %d client(s) ================", clients)
	}

	serverAddress = fmt.Sprintf("%s:%d", serverIP, serverPort)

	log.Println("Starting server...")
	go server()
	time.Sleep(time.Second) // wait for server to come up

	resultChan := make(chan int, clients)
	var wg sync.WaitGroup

	clientConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	}
	if useRes {
		clientConfig.ClientSessionCache = tls.NewLRUClientSessionCache(0) // 0 = default capacity
	}

	log.Printf("Starting %d client(s)...\n", clients)
	for i := 0; i < clients; i++ {

		wg.Add(1)
		go client(clientConfig, resultChan, &wg)
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
	hsps := float64(total) / runTime.Seconds()
	log.Printf("Handshakes per second: %s", strconv.FormatFloat(hsps, 'f', -1, 64))
	println()

	serverPort++

	return total
}
