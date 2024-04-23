package main

import (
	"crypto/tls"
	"log"
	"sync"
	"time"
)

func client(serverAddr string, useRes bool, h chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
	}
	if useRes {
		config.ClientSessionCache = tls.NewLRUClientSessionCache(0) // 0 = default capacity
	}

	handshakes := 0
	startTime := time.Now()
	endTime := startTime.Add(runTime)
	for time.Now().Before(endTime) {
		err := makeRequest(serverAddr, config)
		if err == nil {
			handshakes++
		}
	}
	h <- handshakes
}

func makeRequest(serverAddr string, config *tls.Config) error {
	conn, err := tls.Dial("tcp", serverAddr, config)
	if err != nil {
		log.Printf("client: error dialling: %s", err)
		return err
	}
	defer func(conn *tls.Conn) {
		err = conn.Close()
		if err != nil {
			log.Printf("client: error closing connection: %s", err)
		}
	}(conn)

	err = write(conn, "hello\n")
	if err != nil {
		log.Printf("client: error writing to connection: %s", err)
	}

	err = read(conn)
	if err != nil {
		log.Printf("client: error reading from connection: %s", err)
	}
	return nil
}
