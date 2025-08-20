package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ar3ty/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func replaceProfane(text string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Fields(text)

	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, ok := badWords[lowered]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")

}

func validateChirp(text string) (string, error) {
	const maxChirpLength = 140
	if len(text) > maxChirpLength {
		return "", errors.New("text is too long")
	}

	cleaned := replaceProfane(text)

	return cleaned, nil
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type chirpReq struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	chirpParams := chirpReq{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleanedText, err := validateChirp(chirpParams.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is not valid", err)
		return
	}

	params := database.CreateChirpParams{
		Body:   cleanedText,
		UserID: chirpParams.UserID,
	}

	newChirp, err := cfg.dbQueries.CreateChirp(req.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	response := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
}
