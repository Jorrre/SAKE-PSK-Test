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

const runTime = time.Second * 60
const serverAddr = "127.0.0.1:2208"

func main() {
	log.SetFlags(log.Lmicroseconds)
	parallelClients := []int{1, 2, 4, 6, 8, 10}

	log.Println("Starting server...")
	go server(serverAddr)
	time.Sleep(time.Second) // wait for server to come up

	var fullResults, pskResults []float64
	for _, clients := range parallelClients {
		for _, useRes := range []bool{false, true} {
			result := runTest(clients, useRes)
			if useRes {
				pskResults = append(pskResults, result)
			} else {
				fullResults = append(fullResults, result)
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

func logResult(result []float64) {
	strArr := make([]string, len(result))
	for i, v := range result {
		strArr[i] = fmt.Sprintf("%.1f", v) // You can adjust the formatting as needed
	}
	res := strings.Join(strArr, ", ")
	fmt.Printf("[%s]\n", res)
}
