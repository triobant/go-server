package main

import (
    "encoding/json"
    "net/http"

    "github.com/triobant/go-server/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Password    string  `json:"password"`
        Email       string  `json:"email"`
    }
    type response struct {
        User
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramerters")
        return
    }

    user, err := cfg.DB.GetUserByEmail(params.Email)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
        return
    }

    err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
        return
    }

    respondWithJSON(w, http.StatusOK, response{
        User: User{
            ID:     user.ID,
            Email:  user.Email,
        },
    })
}