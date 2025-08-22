package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ar3ty/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeSubscription(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	gotAPI, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get apikey", err)
		return
	}
	if gotAPI != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Not enough permissions", err)
		return
	}

	reqToParse := request{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&reqToParse)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	if reqToParse.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueries.UpdateUserSubscription(req.Context(), reqToParse.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
