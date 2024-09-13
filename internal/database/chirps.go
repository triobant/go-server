package database

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    id := len(dbStructure.Chirps) + 1
    chirp := Chirp{
        ID:     id,
        Body:   body,
    }
    dbStructure.Chirps[id] = chirp

    err = db.writeDB(dbStructure)
    if err != nil {
        return Chirp{}, err
    }

    return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return nil, err
    }

    chirps := make([]Chirp, 0, len(dbStructure.Chirps))
    for _, chirp := range dbStructure.Chirps {
        chirps = append(chirps, chirp)
    }

    return chirps, nil
}

// GetChirp returns one Chirp by ID
func (db *DB) GetChirp(id int) (Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    chirp, ok := dbStructure.Chirps[id]
    if !ok {
        return Chirp{}, ErrNotExist
    }

    return chirp, nil
}
