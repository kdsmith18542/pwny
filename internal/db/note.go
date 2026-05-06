// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	HostID      string    `json:"host_id,omitempty"`
	Title       string    `json:"title"`
	Content     string    `json:"content,omitempty"`
	NType       string    `json:"ntype"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Database) CreateNote(workspaceID, hostID, title, content, ntype string) (*Note, error) {
	n := &Note{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		HostID:      hostID,
		Title:       title,
		Content:     content,
		NType:       ntype,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := d.Exec(
		`INSERT INTO notes (id, workspace_id, host_id, title, content, ntype, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.ID, n.WorkspaceID, n.HostID, n.Title, n.Content, n.NType, n.CreatedAt, n.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return n, nil
}

func (d *Database) GetNote(id string) (*Note, error) {
	n := &Note{}
	err := d.QueryRow(
		`SELECT id, workspace_id, host_id, title, content, ntype, created_at, updated_at FROM notes WHERE id = ?`, id,
	).Scan(&n.ID, &n.WorkspaceID, &n.HostID, &n.Title, &n.Content, &n.NType, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("note not found: %s", id)
	}
	return n, nil
}

func (d *Database) ListNotes(workspaceID string) ([]Note, error) {
	rows, err := d.Query(
		`SELECT id, workspace_id, host_id, title, content, ntype, created_at, updated_at FROM notes WHERE workspace_id = ? ORDER BY updated_at DESC`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.WorkspaceID, &n.HostID, &n.Title, &n.Content, &n.NType, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, rows.Err()
}

func (d *Database) UpdateNote(id string, updater func(n *Note)) (*Note, error) {
	n, err := d.GetNote(id)
	if err != nil {
		return nil, err
	}

	updater(n)
	n.UpdatedAt = time.Now()

	_, err = d.Exec(
		`UPDATE notes SET title=?, content=?, ntype=?, updated_at=? WHERE id=?`,
		n.Title, n.Content, n.NType, n.UpdatedAt, n.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return n, nil
}

func (d *Database) DeleteNote(id string) error {
	_, err := d.Exec(`DELETE FROM notes WHERE id = ?`, id)
	return err
}
