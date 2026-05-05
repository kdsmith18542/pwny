package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Database) CreateWorkspace(name, description string) (*Workspace, error) {
	w := &Workspace{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := d.Exec(
		`INSERT INTO workspaces (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		w.ID, w.Name, w.Description, w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	return w, nil
}

func (d *Database) GetWorkspace(id string) (*Workspace, error) {
	w := &Workspace{}
	err := d.QueryRow(
		`SELECT id, name, description, created_at, updated_at FROM workspaces WHERE id = ?`, id,
	).Scan(&w.ID, &w.Name, &w.Description, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("workspace not found: %s", id)
	}
	return w, nil
}

func (d *Database) ListWorkspaces() ([]Workspace, error) {
	rows, err := d.Query(`SELECT id, name, description, created_at, updated_at FROM workspaces ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		var w Workspace
		if err := rows.Scan(&w.ID, &w.Name, &w.Description, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, w)
	}

	return workspaces, rows.Err()
}

func (d *Database) DeleteWorkspace(id string) error {
	_, err := d.Exec(`DELETE FROM workspaces WHERE id = ?`, id)
	return err
}
