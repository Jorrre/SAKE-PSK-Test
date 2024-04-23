package main

import (
	"crypto/tls"
	"log"
	"net"
)

const serverCertPublic = "server.crt"
const serverCertPrivate = "server.key"

func server(serverAddr string) {
	cer, err := tls.LoadX509KeyPair(serverCertPublic, serverCertPrivate)
	if err != nil {
		log.Fatalf("server: error reading certificate: %s", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS13,
		ServerName:   serverAddr,
	}
	ln, err := tls.Listen("tcp", serverAddr, config)
	if err != nil {
		log.Fatalf("server: error listening on %s: %s", serverAddr, err)
		return
	}
	log.Printf("Server up and running on %s", serverAddr)

	defer func(ln net.Listener) {
		err = ln.Close()
		if err != nil {
			log.Fatalf("server: error closing listener: %s", err)
		}
	}(ln)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("server: error accepting connection: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("server: error closing connection: %s", err)
		}
	}(conn)

	err := read(conn)
	if err != nil {
		log.Printf("server: error reading from connection: %s", err)
	}

	res := "world\n"
	err = write(conn, res)
	if err != nil {
		log.Printf("server: error writing to connection: %s", err)
	}
}
