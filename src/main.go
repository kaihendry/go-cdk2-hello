package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Set at build time so local, CDK, and fast code deployments report the same revision.
var version = "dev"

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Version", version)
	if r.Method != http.MethodHead {
		if _, err := fmt.Fprint(w, "CICD World ... "+version); err != nil {
			slog.Error("error writing response", "error", err)
		}
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           http.HandlerFunc(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("listening", "address", server.Addr, "version", version)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("error listening", "error", err)
		os.Exit(1)
	}
}
