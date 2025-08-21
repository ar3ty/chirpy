package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ar3ty/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type request struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn int    `json:"expires_in_seconds"`
	}

	reqToParse := request{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqToParse)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if reqToParse.ExpiresIn == 0 || reqToParse.ExpiresIn > 3600 {
		reqToParse.ExpiresIn = 3600
	}

	user, err := cfg.dbQueries.GetUserByEmail(req.Context(), reqToParse.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(reqToParse.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Duration(reqToParse.ExpiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't login", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
