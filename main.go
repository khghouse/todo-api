package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	log.Println("server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
