package main

import (
	"SAKE-TLS_Test/utils"
	"crypto/tls"
	"log"
	"net"
)

const serverCertPublic = "server.crt"
const serverCertPrivate = "server.key"

func main() {
	Server()
}

func Server() {
	cer, err := tls.LoadX509KeyPair(serverCertPublic, serverCertPrivate)
	if err != nil {
		log.Fatalf("server: error reading certificate: %s", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		MinVersion:   tls.VersionTLS13,
	}

	port := "2208"
	ln, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		log.Fatalf("server: error listening on port %d: %s", port, err)
		return
	}
	log.Printf("Server up and running on %s", ln.Addr().String())

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

	err := utils.Read(conn)
	if err != nil {
		log.Printf("server: error reading from connection: %s", err)
	}

	res := "world\n"
	err = utils.Write(conn, res)
	if err != nil {
		log.Printf("server: error writing to connection: %s", err)
	}
}
