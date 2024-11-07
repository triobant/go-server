package main

import (
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "time"

    "github.com/triobant/go-server/internal/auth"
    "github.com/triobant/go-server/internal/database"
    "github.com/google/uuid"
)

type Chirp struct {
    ID          uuid.UUID   `json:"id"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    UserID      uuid.UUID   `json:"user_id"`
    Body        string      `json:"body"`

}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body    string `json:"body"`
    }

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
        return
    }
    userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
        return
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    cleaned, err := validateChirp(params.Body)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err.Error(), err)
        return
    }

    chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
        UserID: userID,
        Body:   cleaned,
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        return
    }

    respondWithJSON(w, http.StatusCreated, Chirp{
        ID:         chirp.ID,
        CreatedAt:  chirp.CreatedAt,
        UpdatedAt:  chirp.UpdatedAt,
        UserID:     chirp.UserID,
        Body:       chirp.Body,
    })
}

func validateChirp(body string) (string, error) {
    const maxChirpLength = 140
    if len(body) > maxChirpLength {
        return "", errors.New("Chirp is too long")
    }

    badWords := map[string]struct{}{
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }
    cleaned := getCleanedBody(body, badWords)
    return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
    words := strings.Split(body, " ")
    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := badWords[loweredWord]; ok {
            words[i] = "****"
        }
    }
    cleaned := strings.Join(words, " ")
    return cleaned
}
