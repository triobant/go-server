package database

import (
    "time"
)

type RefreshToken struct {
    UserID      int         `json:"user_id"`
    Token       string      `json:"token"`
    ExpiresAt   time.Time   `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userID int, token string) error {
    dbStructure, err := db.loadDB()
    if err != nil {
        return err
    }

    refreshToken := RefreshToken{
        UserID:         userID,
        Token:          token,
        ExpiresAt:      time.Now().Add(time.Hour),
    }
    dbStructure.RefreshTokens[token] = refreshToken

    err = db.writeDB(dbStructure)
    if err != nil {
        return err
    }
    return nil
}
