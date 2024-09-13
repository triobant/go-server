package database

import (
    "encoding/json"
    "errors"
    "os"
    "sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps  map[int]Chirp `json:"chirps"`
	Users   map[int]User `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
    db := &DB{
        path:   path,
        mu:     &sync.RWMutex{},
    }
    err := db.ensureDB()
    return db, err
}

func (db *DB) createDB() error {
    dbStructure := DBStructure{
        Chirps: map[int]Chirp{},
    }
    return db.writeDB(dbStructure)
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
    _, err := os.ReadFile(db.path)
    if errors.Is(err, os.ErrNotExist) {
        return db.createDB()
    }
    return err
}

// ResetDB resets the database
func (db *DB) ResetDB() error {
    err := os.Remove(db.path)
    if errors.Is(err, os.ErrNotExist) {
        return nil
    }
    return db.ensureDB()
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
    db.mu.RLock()
    defer db.mu.RUnlock()

    dbStructure := DBStructure{}
    dat, err := os.ReadFile(db.path)
    if errors.Is(err, os.ErrNotExist) {
        return dbStructure, err
    }
    err = json.Unmarshal(dat, &dbStructure)
    if err != nil {
        return dbStructure, err
    }

    return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
    db.mu.Lock()
    defer db.mu.Unlock()

    dat, err := json.Marshal(dbStructure)
    if err != nil {
        return err
    }

    err = os.WriteFile(db.path, dat, 0600)
    if err != nil {
        return err
    }
    return nil
}
