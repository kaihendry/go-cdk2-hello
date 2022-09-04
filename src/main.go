package main

import (
	"net/http"
	"os"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/apex/gateway/v2"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Version", os.Getenv("VERSION"))
		w.Write([]byte("Hallo World ... " + os.Getenv("VERSION")))
	})

	port := os.Getenv("_LAMBDA_SERVER_PORT")

	var err error

	if port == "" {
		log.SetHandler(text.Default)
		err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	} else {
		log.SetHandler(jsonhandler.Default)
		err = gateway.ListenAndServe("", nil)
	}
	log.Fatalf("failed to start server: %v", err)
}
