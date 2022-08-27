package main

import (
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	jsonhandler "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"

	"github.com/apex/gateway/v2"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World ... " + time.Now().String()))
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
