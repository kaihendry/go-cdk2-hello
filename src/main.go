package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/apex/gateway/v2"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		build, ok := debug.ReadBuildInfo()
		if !ok {
			http.Error(w, "No build info available", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Version", build.Main.Version)
		_, err := w.Write([]byte("Hallo World ... " + build.Main.Version))
		if err != nil {
			slog.Error("error writing response", "error", err)
		}
	})

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	var err error

	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		err = gateway.ListenAndServe("", nil)
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	}
	slog.Error("error listening", "error", err)
}
