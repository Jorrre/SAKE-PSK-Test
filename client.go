package main

import (
	"crypto/tls"
	"log"
	"sync"
	"time"
)

func client(serverAddress string, useRes bool, h chan int, wg *sync.WaitGroup) {
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
		err := makeRequest(serverAddress, config)
		if err == nil {
			handshakes++
		}
	}
	h <- handshakes
}

func makeRequest(serverAddress string, config *tls.Config) error {
	conn, err := tls.Dial("tcp", serverAddress, config)
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

	_, err = conn.Write([]byte("hello\n"))
	if err != nil {
		log.Printf("client: error writing to connection: %s", err)
		return err
	}

	buf := make([]byte, 100)
	_, err = conn.Read(buf)
	if err != nil {
		log.Printf("client: error reading from connection: %s", err)
		return err
	}
	return nil
}
