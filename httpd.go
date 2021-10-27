package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/cfg"
)

func serveHTTP(handler http.Handler) {
	server := &http.Server{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  300 * time.Second,
		Handler:      handler,
	}

	if strings.HasPrefix(cfg.ListenAddr, "/") {
		startUnixSocketServer(server, cfg.ListenAddr)
	} else {
		server.Addr = cfg.ListenAddr
		startHTTPServer(server)
	}
}

func startUnixSocketServer(server *http.Server, socketFile string) {
	err := os.Remove(socketFile)
	if err != nil {
		log.Println("Failed to remove the socket file", socketFile)
	}

	listener, err := net.Listen("unix", socketFile)
	if err != nil {
		log.Fatalf("Failed to start a server: %v", err)
	}
	defer listener.Close()

	if err := os.Chmod(socketFile, 0666); err != nil {
		log.Fatalf("Failed to set socket permissions: %v", err)
	}

	log.Printf("Listening on Unix socket %s", cfg.ListenAddr)
	if err := server.Serve(listener); err != http.ErrServerClosed {
		log.Fatalf("Failed to start a server: %v", err)
	}
}

func startHTTPServer(server *http.Server) {
	log.Printf("Listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to start a server: %v", err)
	}
}
