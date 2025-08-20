package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorJSON struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with %d error: %s\n", code, msg)
	}
	errorToSend := errorJSON{
		Error: msg + ":" + err.Error(),
	}
	respondWithJSON(w, code, errorToSend)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
