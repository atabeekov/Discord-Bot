package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database struct wraps the SQLite connection
type Database struct {
	DB *sql.DB
}

// InitDB initializes the SQLite database and creates the necessary table
func InitDB(dbPath string) (*Database, error) {
	log.Println("Initializing database...")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Failed to open database: %v\n", err)
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Создаем таблицу, если она еще не существует
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		language TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Failed to create users table: %v\n", err)
		return nil, fmt.Errorf("failed to create users table: %v", err)
	}

	// Проверяем структуру таблицы
	rows, err := db.Query(`PRAGMA table_info(users);`)
	if err != nil {
		log.Printf("Failed to query table structure: %v\n", err)
		return nil, fmt.Errorf("failed to query table structure: %v", err)
	}
	defer rows.Close()

	// Проверяем наличие колонки "id"
	columnExists := false
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull, pk int
		var dfltValue interface{} // default value может быть NULL

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &dfltValue, &pk); err != nil {
			log.Printf("Failed to scan table structure row: %v\n", err)
			return nil, fmt.Errorf("failed to scan table structure row: %v", err)
		}

		if name == "id" { // Проверяем корректное имя колонки
			columnExists = true
			break
		}
	}

	// Если колонки "id" нет, выводим предупреждение
	if !columnExists {
		log.Println("Warning: 'id' column is missing in the users table. Please check your database schema.")
	}

	log.Println("Database initialized successfully.")
	return &Database{DB: db}, nil
}

// GetUserLanguage retrieves the preferred language for a user
func (db *Database) GetUserLanguage(userID string) (string, error) {
	log.Printf("Retrieving language for user %s...\n", userID)

	var language string
	query := `SELECT language FROM users WHERE id = ?`
	err := db.DB.QueryRow(query, userID).Scan(&language)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User %s not found. Returning default language 'en'.\n", userID)
			return "en", nil // Default language
		}
		log.Printf("Failed to retrieve language for user %s: %v\n", userID, err)
		return "", fmt.Errorf("failed to get language: %v", err)
	}

	log.Printf("Language for user %s is %s.\n", userID, language)
	return language, nil
}

// SetUserLanguage stores or updates the preferred language for a user
func (db *Database) SetUserLanguage(userID, language string) error {
	log.Printf("Setting language for user %s to %s...\n", userID, language)

	query := `
	INSERT INTO users (id, language)
	VALUES (?, ?)
	ON CONFLICT(id) DO UPDATE SET language = excluded.language;`
	_, err := db.DB.Exec(query, userID, language)
	if err != nil {
		log.Printf("Failed to set language for user %s: %v\n", userID, err)
		return fmt.Errorf("failed to set language: %v", err)
	}

	log.Printf("Language set successfully for user %s.\n", userID)
	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	log.Println("Closing database connection...")
	err := db.DB.Close()
	if err != nil {
		log.Printf("Failed to close database: %v\n", err)
		return err
	}
	log.Println("Database connection closed successfully.")
	return nil
}
