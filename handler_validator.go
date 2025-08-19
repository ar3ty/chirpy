package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type cleaned struct {
	CleanedBody string `json:"cleaned_body"`
}

type ValidJSON struct {
	Valid bool `json:"valid"`
}

func replaceProfane(text string) string {
	badWords := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}
	words := strings.Fields(text)
	for badWord, replacevalue := range badWords {
		for i, word := range words {
			if strings.ToLower(word) == badWord {
				words[i] = replacevalue
			}
		}
	}

	return strings.Join(words, " ")

}

func handlerChirpValidator(w http.ResponseWriter, req *http.Request) {
	params := parameters{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respond := cleaned{
		CleanedBody: replaceProfane(params.Body),
	}

	respondWithJSON(w, http.StatusOK, respond)
}
