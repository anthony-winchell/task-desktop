package main

import (
	"fmt"
	"strings"
)

func (a *App) GetTasks() ([]Task, error) {
	db, err := a.database()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
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
	db, err := a.database()
	if err != nil {
		return err
	}

	title = strings.TrimSpace(title)

	if title == "" {
		return fmt.Errorf("title is required")
	}

	_, err = db.Exec(
		"INSERT INTO tasks (title) values (?)",
		title,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	return nil
}

func (a *App) DeleteTask(id int) error {
	db, err := a.database()
	if err != nil {
		return err
	}

	result, err := db.Exec(
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
	db, err := a.database()
	if err != nil {
		return err
	}

	result, err := db.Exec(
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

func (a *App) UpdateTitle(id int, newTitle string) error {
	db, err := a.database()
	if err != nil {
		return err
	}

	trimmed := strings.TrimSpace(newTitle)

	if trimmed == "" {
		return fmt.Errorf("title is required")
	}

	result, err := db.Exec(
		"UPDATE tasks SET title = ? WHERE id = ?",
		trimmed,
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
