package main

import (
	"fmt"
	"net/http"
)

func main() {
	var log = logging.MustGetLogger("example")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world")
	})
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API is up")
	})
	http.ListenAndServe(":8080", nil)
}
