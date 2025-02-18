package models

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"time"

	// _ "github.com/lib/pq"
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var Db *sqlx.DB // *sql.DB

func InitDB() {
	var err error
	Db, err = sqlx.Open("sqlite", ":memory:") // sql.Open("postgres", "dbname=chitchat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// create the threads table
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
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return uuid
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}

func InsertInitialDB() {
	var err error

	// Insert initial data
	_, err = Db.Exec(`
        INSERT INTO users (uuid, name, email, password, created_at) VALUES (?, ?, ?, ?, ?)
    `, createUUID(), "taewony", "taewony@gmail.com", Encrypt("password123"), time.Now())
	if err != nil {
		log.Fatal(err)
	}
	_, err = Db.Exec(`
        INSERT INTO sessions (uuid, email, user_id, created_at) VALUES (?, ?, ?, ?)
    `, createUUID(), "taewony@gmail.com", 1, time.Now()) // Assuming user_id is 1
	if err != nil {
		log.Fatal(err)
	}
	_, err = Db.Exec(`
        INSERT INTO threads (uuid, topic, user_id, created_at) VALUES (?, ?, ?, ?)
    `, createUUID(), "Sample Thread", 1, time.Now()) // Assuming user_id is 1
	if err != nil {
		log.Fatal(err)
	}
	_, err = Db.Exec(`
        INSERT INTO posts (uuid, body, user_id, thread_id, created_at) VALUES (?, ?, ?, ?, ?)
    `, createUUID(), "This is a sample post body.", 1, 1, time.Now()) // Assuming user_id is 1 and thread_id is 1
	if err != nil {
		log.Fatal(err)
	}
}
