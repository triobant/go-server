package database

import "errors"

type User struct {
	ID              int     `json:"id"`
	Email           string  `json:"email"`
    HashedPassword  string  `json:"hashed_password"`
}

var ErrAlreadyExists = errors.New("already exists")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
    if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
        return User{}, ErrAlreadyExists
    }

    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    id := len(dbStructure.Users) + 1
    user := User{
        ID:             id,
        Email:          email,
        HashedPassword: hashedPassword,
    }
    dbStructure.Users[id] = user

    err = db.writeDB(dbStructure)
    if err != nil {
        return User{}, err
    }

    return user, nil
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

// GetUser returns one User by ID
func (db *DB) GetUser(id int) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    user, ok := dbStructure.Users[id]
    if !ok {
        return User{}, ErrNotExist
    }

    return user, nil
}

// GetUserByEmail returns one User by email
func (db *DB) GetUserByEmail(email string) (User, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return User{}, err
    }

    for _, user := range dbStructure.Users {
        if user.Email == email {
            return user, nil
        }
    }

    return User{}, ErrNotExist
}