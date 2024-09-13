package main

import (
    "encoding/json"
    "net/http"
)

type User struct {
    ID      int     `json:"id"`
    Body    string  `json:"email"`

}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

    user, err := cfg.DB.CreateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
        return
    }

    respondWithJSON(w, http.StatusCreated, User{
        ID:     user.ID,
        Body:   user.Body,
    })
}
