package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFolder = ".tasks"
	dbName   = "tasks.db"
)

var DB *Database

type Database struct {
	db *sql.DB
}

func getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	dbPath := filepath.Join(homeDir, dbFolder, dbName)

	err = os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	return dbPath
}

func Get() *Database {
	if DB == nil {
		connect()
	}
	return DB
}

func connect() {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	DB = &Database{
		db: db,
	}

	// TODO run migrations
}
