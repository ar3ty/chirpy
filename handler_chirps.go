package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ar3ty/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	chirpParams := chirpReq{}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't authorize the post", err)
		return
	}

	jwtID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate the session", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&chirpParams)
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
		UserID: jwtID,
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

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	var err error
	var chirps []database.Chirp
	order := "asc"

	sortingOrderStr := req.URL.Query().Get("sort")
	if sortingOrderStr == "desc" {
		order = sortingOrderStr
	}

	if author := req.URL.Query().Get("author_id"); author != "" {
		authorID, parseErr := uuid.Parse(author)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse id", parseErr)
			return
		}
		chirps, err = cfg.dbQueries.GetChirpsByAuthorID(req.Context(), authorID)
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(req.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		if order == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		} else {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		}
	})

	response := []Chirp{}
	for _, item := range chirps {
		response = append(response, Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID:    item.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpIdStr := req.PathValue("chirpID")
	id, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(req.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Couldn't find a chirp by ID", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get a chirp by ID", err)
		return
	}

	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	chirpIDString := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Couldn't find a chirp by ID", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't get a chirp by ID", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not enough rights to delete", err)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
