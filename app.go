package main

import (
	"context"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// App struct
type App struct {
	ctx   context.Context
	db    *sql.DB
	dbErr error
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.dbErr = a.initDB()
	if a.dbErr != nil {
		log.Printf("initialize database: %v", a.dbErr)
	}
}

func (a *App) shutdown(_ context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Close database: %v", err)
		}
	}
}
