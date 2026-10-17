package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		if err = db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}
		return fmt.Errorf("connect to database: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            completed INTEGER NOT NULL DEFAULT 0
        )
    `)
	if err != nil {
		if err = db.Close(); err != nil {
			return fmt.Errorf("close database: %w", err)
		}
		return fmt.Errorf("create tasks table: %w", err)
	}

	a.db = db
	log.Printf("Database initialized at %s", dbPath)

	if err := a.seedDB(); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	return nil
}

func (a *App) seedDB() error {
	var count int

	err := a.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return fmt.Errorf("count tasks: %w", err)
	}

	if count == 0 {
		_, err = a.db.Exec(`
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.initDB(); err != nil {
		log.Printf("Initialize database: %v", err)
	}
}

func (a *App) shutdown(_ context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Close database: %v", err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetTasks() ([]Task, error) {
	rows, err := a.db.Query(
		"SELECT id, title, completed FROM tasks ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}

	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task
		var completed int
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&completed,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		task.Completed = completed != 0

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (a *App) AddTask(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return fmt.Errorf("title is required")
	}

	_, err := a.db.Exec(
		"INSERT INTO tasks (title) values (?)",
		title,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	return nil
}

func (a *App) DeleteTask(id int) error {
	result, err := a.db.Exec(
		"DELETE FROM tasks WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %d not found", id)
	}

	return nil
}

func (a *App) SetTaskCompleted(id int, completed bool) error {
	result, err := a.db.Exec(
		"UPDATE tasks SET completed = ? WHERE id = ?",
		completed,
		id,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task %d not found", id)
	}

	return nil
}
