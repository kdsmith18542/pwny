// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"fmt"
	"time"
)

type EventLog struct {
	ID          int       `json:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	Level       string    `json:"level"`
	Module      string    `json:"module,omitempty"`
	Message     string    `json:"message"`
	Detail      string    `json:"detail,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (d *Database) CreateEvent(workspaceID, level, module, message, detail string) (*EventLog, error) {
	e := &EventLog{
		WorkspaceID: workspaceID,
		Level:       level,
		Module:      module,
		Message:     message,
		Detail:      detail,
		CreatedAt:   time.Now(),
	}

	res, err := d.Exec(
		`INSERT INTO event_log (workspace_id, level, module, message, detail, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		e.WorkspaceID, e.Level, e.Module, e.Message, e.Detail, e.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create event log: %w", err)
	}

	id, _ := res.LastInsertId()
	e.ID = int(id)

	return e, nil
}

func (d *Database) ListEvents(workspaceID string, limit int) ([]EventLog, error) {
	query := `SELECT id, workspace_id, level, module, message, COALESCE(detail,''), created_at FROM event_log`
	args := []interface{}{}

	if workspaceID != "" {
		query += ` WHERE workspace_id = ?`
		args = append(args, workspaceID)
	}

	query += ` ORDER BY created_at DESC`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := d.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventLog
	for rows.Next() {
		var e EventLog
		if err := rows.Scan(&e.ID, &e.WorkspaceID, &e.Level, &e.Module, &e.Message, &e.Detail, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}
