package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mabzd/snorlax/internal/config"
)

func InitDB(cfg config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	var db *sql.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 3)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	log.Println("Successfully connected to database")

	if err := createMigrationsTable(db); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %v", err)
	}

	return db, nil
}

func GetAppliedMigrationFileNamesSet(db *sql.DB) (map[string]struct{}, error) {
	query := "SELECT file_name FROM migrations"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fileNames := make(map[string]struct{})
	for rows.Next() {
		var fileName string
		if err := rows.Scan(&fileName); err != nil {
			return nil, err
		}
		fileNames[fileName] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fileNames, nil
}

func MustApplyMigration(db *sql.DB, sql string, fileName string) {
	start := time.Now()
	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("Failed to apply migration '%s': %v", fileName, err)
	}

	insertSql := `
		INSERT INTO migrations (file_name, elapsed_ms)
		VALUES ($1, $2)
	`
	elapsedMs := time.Since(start).Milliseconds()
	if _, err := db.Exec(insertSql, fileName, elapsedMs); err != nil {
		log.Fatalf("Failed to insert migration record for '%s': %v", fileName, err)
	}
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			file_name VARCHAR(255) NOT NULL,
			elapsed_ms BIGSERIAL NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}
