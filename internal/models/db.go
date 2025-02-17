package models

import (
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

func InitDB() {
	var err error
	Db, err = sqlx.Open("sqlite", ":memory:") // Change "sqlite3" to "sqlite"
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	_, err = Db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid       VARCHAR(64) NOT NULL UNIQUE,
            name       VARCHAR(255),
            email      VARCHAR(255) NOT NULL UNIQUE,
            password   VARCHAR(255) NOT NULL,
            created_at TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid       VARCHAR(64) NOT NULL UNIQUE,
            email      VARCHAR(255),
            user_id    INTEGER REFERENCES users(id),
            created_at TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        CREATE TABLE IF NOT EXISTS threads (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid       VARCHAR(64) NOT NULL UNIQUE,
            topic      TEXT,
            user_id    INTEGER REFERENCES users(id),
            created_at TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            uuid       VARCHAR(64) NOT NULL UNIQUE,
            body       TEXT,
            user_id    INTEGER REFERENCES users(id),
            thread_id  INTEGER REFERENCES threads(id),
            created_at TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		log.Fatal(err)
	}

	// Insert initial data
	_, err = Db.Exec(`
        INSERT INTO users (uuid, name, email, password, created_at) VALUES (?, ?, ?, ?, ?)
    `, CreateUUID(), "taewony", "taewony@gmail.com", Encrypt("password123"), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        INSERT INTO sessions (uuid, email, user_id, created_at) VALUES (?, ?, ?, ?)
    `, CreateUUID(), "taewony@gmail.com", 1, time.Now()) // Assuming user_id is 1
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        INSERT INTO threads (uuid, topic, user_id, created_at) VALUES (?, ?, ?, ?)
    `, CreateUUID(), "Sample Thread", 1, time.Now()) // Assuming user_id is 1
	if err != nil {
		log.Fatal(err)
	}

	_, err = Db.Exec(`
        INSERT INTO posts (uuid, body, user_id, thread_id, created_at) VALUES (?, ?, ?, ?, ?)
    `, CreateUUID(), "This is a sample post body.", 1, 1, time.Now()) // Assuming user_id is 1 and thread_id is 1
	if err != nil {
		log.Fatal(err)
	}
}

// Create a random UUID
func CreateUUID() string {
	u := uuid.New()
	return u.String()
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

// Check if the password is correct
func CheckPassword(email, password string) (valid bool, err error) {
	user, err := UserByEmail(email)
	if err != nil {
		return
	}
	return Encrypt(password) == user.Password, nil
}
