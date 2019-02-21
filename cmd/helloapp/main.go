package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const name = "helloapp"

var hostname string

func sayHello(w http.ResponseWriter, r *http.Request) {
	// Log on each request and say hello
	log.Print("Request: ", r.RequestURI, " From: ", r.RemoteAddr)
	_, err := fmt.Fprintf(w, "hello %s, this is %s", r.RemoteAddr, hostname)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func main() {
	var (
		err  error
		port uint
		ctx  = context.Background()
	)
	log.Print(name, ": starting")

	// CLI flags
	flag.UintVar(&port, "port", 8080, "Port for service to listen on")
	flag.Parse()

	// Setup channels to capture SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 2)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Setup the HTTP server
	http.HandleFunc("/", sayHello)
	server := &http.Server{}

	// Shutdown the server gracefully
	go func(ctx context.Context) {
		<-sigChan
		log.Printf("%s: shutting down", name)

		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		shutdownErr := server.Shutdown(ctx)
		if err != nil {
			log.Fatal("Error shutting down:", shutdownErr)
		}
	}(ctx)

	// Discover the hostname
	hostname, err = os.Hostname()
	if err != nil {
		log.Print("Could not find hostname: ", err)
		hostname = "unknown"
	}
	log.Print(name, ": my hostname is: ", hostname)

	// Start the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port:", err)
	}
	defer listener.Close()

	// Start the server
	log.Print(name, ": listening on port: ", listener.Addr().(*net.TCPAddr).Port)
	err = server.Serve(listener)
	if err != nil && context.Canceled == nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
