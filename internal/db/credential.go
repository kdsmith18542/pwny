// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	HostID      string    `json:"host_id,omitempty"`
	Username    string    `json:"username"`
	PasswordEnc string    `json:"-"`
	HashEnc     string    `json:"-"`
	Type        string    `json:"type"`
	Module      string    `json:"module,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (d *Database) CreateCredential(workspaceID, hostID, username, passwordEnc, hashEnc, credType, module string) (*Credential, error) {
	c := &Credential{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		HostID:      hostID,
		Username:    username,
		PasswordEnc: passwordEnc,
		HashEnc:     hashEnc,
		Type:        credType,
		Module:      module,
		CreatedAt:   time.Now(),
	}

	_, err := d.Exec(
		`INSERT INTO credentials (id, workspace_id, host_id, username, password_enc, hash_enc, type, module, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID, c.WorkspaceID, c.HostID, c.Username, c.PasswordEnc, c.HashEnc, c.Type, c.Module, c.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return c, nil
}

func (d *Database) GetCredential(id string) (*Credential, error) {
	c := &Credential{}
	err := d.QueryRow(
		`SELECT id, workspace_id, host_id, username, password_enc, hash_enc, type, module, created_at FROM credentials WHERE id = ?`, id,
	).Scan(&c.ID, &c.WorkspaceID, &c.HostID, &c.Username, &c.PasswordEnc, &c.HashEnc, &c.Type, &c.Module, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("credential not found: %s", id)
	}
	return c, nil
}

func (d *Database) ListCredentials(workspaceID string) ([]Credential, error) {
	rows, err := d.Query(
		`SELECT id, workspace_id, host_id, username, password_enc, hash_enc, type, module, created_at FROM credentials WHERE workspace_id = ? ORDER BY created_at DESC`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []Credential
	for rows.Next() {
		var c Credential
		if err := rows.Scan(&c.ID, &c.WorkspaceID, &c.HostID, &c.Username, &c.PasswordEnc, &c.HashEnc, &c.Type, &c.Module, &c.CreatedAt); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}

	return creds, rows.Err()
}

func (d *Database) DeleteCredential(id string) error {
	_, err := d.Exec(`DELETE FROM credentials WHERE id = ?`, id)
	return err
}
