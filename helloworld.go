package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var (
		dbUser    = os.Getenv("DB_USER") // e.g. 'my-db-user'
		dbPwd     = os.Getenv("DB_PASS") // e.g. 'my-db-password'
		dbTCPHost = os.Getenv("DB_HOST") // e.g. '127.0.0.1' ('172.17.0.1' if deployed to GAE Flex)
		dbPort    = os.Getenv("DB_PORT") // e.g. '3306'
		dbName    = os.Getenv("DB_NAME") // e.g. 'my-database'
	)

	if dbUser == "" {
		log.Fatal("DB_USER environment variable must be set")
	}

	if dbPwd == "" {
		log.Fatal("DB_PASS environment variable must be set")
	}

	if dbTCPHost == "" {
		log.Fatal("DB_HOST environment variable must be set")
	}

	if dbPort == "" {
		log.Fatal("DB_PORT environment variable must be set")
	}

	if dbName == "" {
		log.Fatal("DB_NAME environment variable must be set")
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		fmt.Printf("sql.Open: %v", err)
		return
	}

	results, err := dbPool.Query("SELECT greeting FROM hello")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var herro string

	for results.Next() {

		// for each row, scan the result into our tag composite object
		err = results.Scan(&herro)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(herro)
	}

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
			go handleConnection(conn, herro)
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

func handleConnection(conn net.Conn, herro string) {
	defer conn.Close()

	hello := fmt.Sprintf("%s ", herro)

	for {
		time.Sleep(1 * time.Second)

		var elloh = []string{}

		for i := 0; i < len(hello)-1; i++ {
			elloh = append(elloh, string(hello[i+1]))
		}

		elloh = append(elloh, string(hello[0]))

		hello = strings.Join(elloh, "")

		n, err := conn.Write([]byte(hello))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
