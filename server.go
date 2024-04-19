package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"time"
)

const serverCertPublic = "server.crt"
const serverCertPrivate = "server.key"

func server() {
	cer, err := tls.LoadX509KeyPair(serverCertPublic, serverCertPrivate)
	if err != nil {
		log.Fatalln("server: error reading certificate")
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS13,
		ServerName:   serverAddress,
	}
	ln, err := tls.Listen("tcp", serverAddress, config)
	if err != nil {
		log.Fatalf("server: error listening on %s: %s", serverAddress, err)
		return
	}
	log.Printf("Server up and running on %s", serverAddress)

	defer func(ln net.Listener) {
		err = ln.Close()
		if err != nil {
			log.Fatalln("server: error closing listener")
		}
		time.Sleep(time.Second)
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("server: error accepting connection")
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("server: error closing connection")
		}
	}(conn)

	r := bufio.NewReader(conn)
	for {
		_, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("server: error reading from connection: %s", err)
			}
			return
		}

		_, err = conn.Write([]byte("world!\n"))
		if err != nil {
			log.Println("server: error writing to connection")
			return
		}
	}
}
