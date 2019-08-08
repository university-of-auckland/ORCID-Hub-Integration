//+build !test,heroku

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT not set")
	}
	addr := ":" + port
	http.HandleFunc("/v1/enqueue", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var e Event
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&e)
		if err != nil {
			rw.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprintf(rw, `{"error": %q}`, err.Error())
			return
		}
		e.handle()
	})
	log.Infof("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
