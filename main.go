package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxClients = 10
const runTime = time.Second * 5
const serverIP = "127.0.0.1"

var serverPort = 2208

func main() {
	log.SetFlags(log.Lmicroseconds)
	var fullResults, pskResults [maxClients]float64
	for clients := 1; clients <= maxClients; clients++ {
		for _, useRes := range []bool{true, true} {
			result := runTest(clients, useRes)
			if useRes {
				pskResults[clients-1] = result
			} else {
				fullResults[clients-1] = result
			}
		}
	}
	log.Printf("**************** Results for full handshake, %d seconds ****************", int(runTime.Seconds()))
	logResult(fullResults)
	log.Printf("**************** Results for PSK handshake %d seconds ****************", int(runTime.Seconds()))
	logResult(pskResults)
}

func runTest(clients int, useRes bool) float64 {
	if useRes {
		log.Printf("================ Running TLS-PSK test with %d client(s) ================", clients)
	} else {
		log.Printf("================ Running full TLS handshake test with %d client(s) ================", clients)
	}

	serverAddr := fmt.Sprintf("%s:%d", serverIP, serverPort)

	log.Println("Starting server...")
	go server(serverAddr)
	time.Sleep(time.Second) // wait for server to come up

	resultChan := make(chan int, clients)
	var wg sync.WaitGroup

	log.Printf("Starting %d client(s)...\n", clients)
	for i := 0; i < clients; i++ {
		wg.Add(1)
		go client(serverAddr, useRes, resultChan, &wg)
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

	return hsps
}

func read(conn net.Conn) error {
	r := bufio.NewReader(conn)
	_, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	return nil
}

func write(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func logResult(result [maxClients]float64) {
	strArr := make([]string, len(result))
	for i, v := range result {
		strArr[i] = fmt.Sprintf("%.1f", v) // You can adjust the formatting as needed
	}
	res := strings.Join(strArr, ", ")
	fmt.Printf("[%s]\n", res)
}
