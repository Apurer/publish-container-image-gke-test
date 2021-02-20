package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8443"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	tlsCert, tlsKey := os.Getenv("TLS_CERT"), os.Getenv("TLS_KEY")
	if tlsCert == "" {
		log.Fatal("TLS_CERT environment variable must be set")
	}
	if tlsKey == "" {
		log.Fatal("TLS_KEY environment variable must be set")
	}

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/", hello)

	// start the web server on port and accept requests
	log.Printf("tls cert: %s", tlsCert)
	log.Printf("tls key: %s", tlsKey)
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServeTLS(":"+port, tlsCert, tlsKey, server)
	log.Fatal(err)
}

// hello responds to the request with a plain-text "Hello, world" message.
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Protocol: %s!\n", r.Proto)
	fmt.Fprintf(w, "Hostname: %s\n", host)
	if headerIP := r.Header.Get("X-Forwarded-For"); headerIP != "" {
		fmt.Fprintf(w, "Client IP (X-Forwarded-For): %s\n", headerIP)
	}
}
