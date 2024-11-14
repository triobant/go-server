package main

import (
    "database/sql"
    "encoding/json"
    "errors"
    "net/http"

    "github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Event   string  `json:"event"`
        Data    struct  {
            UserID  uuid.UUID   `json:"user_id"`
        }
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    }

    if params.Event != "user.upgraded" {
        w.WriteHeader(http.StatusNoContent)
        return
    }

    _, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
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
