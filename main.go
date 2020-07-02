package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	})
	fmt.Println("Starting server on the port 8080...")

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}
