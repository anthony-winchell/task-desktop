package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (a *App) database() (*sql.DB, error) {
	if a.dbErr != nil {
		return nil, fmt.Errorf("database initialization failed: %w", a.dbErr)
	}

	if a.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	return a.db, nil
}

func (a *App) initDB() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get config directory: %w", err)
	}

	appDir := filepath.Join(configDir, "TaskTracker")

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return fmt.Errorf("create app directory: %w", err)
	}

	dbPath := filepath.Join(appDir, "tasks.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	initialized := false
	defer func() {
		if !initialized {
			db.Close()
		}
	}()

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            completed INTEGER NOT NULL DEFAULT 0
        )
    `)
	if err != nil {
		return fmt.Errorf("create tasks table: %w", err)
	}

	if err := seedTasks(db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	a.db = db
	initialized = true
	log.Printf("Database initialized at %s", dbPath)

	return nil
}

func seedTasks(db *sql.DB) error {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return fmt.Errorf("count tasks: %w", err)
	}

	if count == 0 {
		_, err = db.Exec(`
				INSERT INTO tasks (title) VALUES
				('Learn Wails'),
				('Build a task tracker')
			`)
		if err != nil {
			return fmt.Errorf("insert tasks: %w", err)
		}
	}

	return nil
}
