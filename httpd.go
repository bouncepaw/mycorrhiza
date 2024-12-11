package main

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
)

func serveHTTP(handler http.Handler) (err error) {
	server := &http.Server{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		IdleTimeout:  300 * time.Second,
		Handler:      handler,
	}

	if strings.HasPrefix(cfg.ListenAddr, "/") {
		err = startUnixSocketServer(server, cfg.ListenAddr)
	} else {
		server.Addr = cfg.ListenAddr
		err = startHTTPServer(server)
	}
	return err
}

func startUnixSocketServer(server *http.Server, socketPath string) error {
	err := os.Remove(socketPath)
	if err != nil {
		slog.Warn("Failed to clean up old socket", "err", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		slog.Error("Failed to start the server", "err", err)
		return err
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	if err := os.Chmod(socketPath, 0666); err != nil {
		slog.Error("Failed to set socket permissions", "err", err)
		return err
	}

	slog.Info("Listening Unix socket", "addr", socketPath)

	if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start the server", "err", err)
		return err
	}

	return nil
}

func startHTTPServer(server *http.Server) error {
	slog.Info("Listening over HTTP", "addr", server.Addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start the server", "err", err)
		return err
	}

	return nil
}
