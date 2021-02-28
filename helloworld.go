package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	http2ProxyPort := os.Getenv("HTTP2_SERVER_PORT")
	if http2ProxyPort == "" {
		log.Fatal("HTTP2_SERVER_PORT environment variable must be set")
	}
	tcpServerPort := os.Getenv("TCP_SERVER_PORT")
	if tcpServerPort == "" {
		log.Fatal("TCP_SERVER_PORT environment variable must be set")
	}

	tlsCert, tlsKey := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")
	if tlsCert == "" {
		log.Fatal("TLS_CERT environment variable must be set")
	}
	if tlsKey == "" {
		log.Fatal("TLS_KEY environment variable must be set")
	}

	cer, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":"+tcpServerPort, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/", hello)

	// start the web server on port and accept requests
	log.Printf("tls cert: %s", tlsCert)
	log.Printf("tls key: %s", tlsKey)
	log.Printf("Server listening on port %s", http2ProxyPort)
	err = http.ListenAndServeTLS(":"+http2ProxyPort, tlsCert, tlsKey, server)
	log.Fatal(err)
}

// hello responds to the request with a plain-text "Hello, world" message.
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Scheme: %s\n", r.URL.Scheme)
	fmt.Fprintf(w, "Protocol: %s!\n", r.Proto)
	fmt.Fprintf(w, "Hostname: %s\n", host)

	// ipAddr := strings.Split(r.RemoteAddr, ":")[0]
	// if ipAddr != "" {
	// 	fmt.Fprintf(w, "Client IP: %s\n", ipAddr)
	// }

	if headerIP := r.Header.Get("X-Forwarded-For"); headerIP != "" {
		fmt.Fprintf(w, "Client IP (X-Forwarded-For): %s\n", headerIP)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		time.Sleep(1 * time.Second)

		n, err := conn.Write([]byte("hello there\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
