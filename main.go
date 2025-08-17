package main

import "net/http"

func main() {
	port := "8080"

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}
