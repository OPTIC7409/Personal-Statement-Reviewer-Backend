package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func InitializeDatabase() bool {
	dbPath := filepath.Join(".", "psr_database.db")
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return false
	}

	if err = db.Ping(); err != nil {
		log.Printf("Error connecting to database: %v", err)
		return false
	}

	fmt.Println("Connected to the local SQLite database")
	return true
}

func CreateTables() bool {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS personal_statements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS feedback (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			statement_id INTEGER NOT NULL,
			feedback_text TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS revisions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			statement_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS user_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			bio TEXT,
			profile_picture_url TEXT,
			preferences TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			action TEXT NOT NULL,
			details TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS ai_detection (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			statement_id INTEGER NOT NULL,
			overall_ai_probability FLOAT NOT NULL,
			flagged_sections TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (statement_id) REFERENCES personal_statements (id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating table: %v", err)
			return false
		}
	}

	fmt.Println("Created tables")
	return true
}

func GetConnection() *sql.DB {
	return db
}

func CloseDatabase() {
	if db != nil {
		db.Close()
	}
}
