//+build !test,heroku !test,container !test,standalone

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT not set")
	}
	addr := ":" + port
	http.HandleFunc("/handle", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var e Event
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&e)
		if err != nil {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprintf(rw, `{"error": %q}`, err.Error())
			return
		}
		msg, err := e.handle()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, `{"error": %q}`, err.Error())
		}
		if msg != "" {
			fmt.Fprintf(rw, `{"message": %q}`, msg)
		} else {
			rw.WriteHeader(http.StatusNoContent)
		}
	})
	log.Infof("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
