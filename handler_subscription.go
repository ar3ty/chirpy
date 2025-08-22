package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeSubscription(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	reqToParse := request{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqToParse)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}
	if reqToParse.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(reqToParse.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user_id", err)
		return
	}

	_, err = cfg.dbQueries.UpdateUserSubscription(req.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't update user subscription", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
