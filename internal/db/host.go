// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Host struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Address     string    `json:"address"`
	MAC         string    `json:"mac,omitempty"`
	OSName      string    `json:"os_name,omitempty"`
	OSFlavor    string    `json:"os_flavor,omitempty"`
	OSSP        string    `json:"os_sp,omitempty"`
	Arch        string    `json:"arch,omitempty"`
	Purpose     string    `json:"purpose,omitempty"`
	Info        string    `json:"info,omitempty"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const hostCols = `id, workspace_id, address, COALESCE(mac,''), COALESCE(os_name,''), COALESCE(os_flavor,''), COALESCE(os_sp,''), COALESCE(arch,''), COALESCE(purpose,''), COALESCE(info,''), state, created_at, updated_at`

func (d *Database) CreateHost(workspaceID, address string) (*Host, error) {
	h := &Host{
		ID:          uuid.New().String(),
		WorkspaceID: workspaceID,
		Address:     address,
		State:       "alive",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := d.Exec(
		`INSERT INTO hosts (id, workspace_id, address, state, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		h.ID, h.WorkspaceID, h.Address, h.State, h.CreatedAt, h.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	return h, nil
}

func (d *Database) GetHost(id string) (*Host, error) {
	h := &Host{}
	err := d.QueryRow(
		`SELECT `+hostCols+` FROM hosts WHERE id = ?`, id,
	).Scan(&h.ID, &h.WorkspaceID, &h.Address, &h.MAC, &h.OSName, &h.OSFlavor, &h.OSSP, &h.Arch, &h.Purpose, &h.Info, &h.State, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("host not found: %s", id)
	}
	return h, nil
}

func (d *Database) ListHosts(workspaceID string) ([]Host, error) {
	rows, err := d.Query(
		`SELECT `+hostCols+` FROM hosts WHERE workspace_id = ? ORDER BY address`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []Host
	for rows.Next() {
		var h Host
		if err := rows.Scan(&h.ID, &h.WorkspaceID, &h.Address, &h.MAC, &h.OSName, &h.OSFlavor, &h.OSSP, &h.Arch, &h.Purpose, &h.Info, &h.State, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		hosts = append(hosts, h)
	}

	return hosts, rows.Err()
}

func (d *Database) DeleteHost(id string) error {
	_, err := d.Exec(`DELETE FROM hosts WHERE id = ?`, id)
	return err
}

func (d *Database) UpdateHost(id string, updater func(h *Host)) (*Host, error) {
	h, err := d.GetHost(id)
	if err != nil {
		return nil, err
	}

	updater(h)
	h.UpdatedAt = time.Now()

	_, err = d.Exec(
		`UPDATE hosts SET mac=?, os_name=?, os_flavor=?, os_sp=?, arch=?, purpose=?, info=?, state=?, updated_at=? WHERE id=?`,
		h.MAC, h.OSName, h.OSFlavor, h.OSSP, h.Arch, h.Purpose, h.Info, h.State, h.UpdatedAt, h.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update host: %w", err)
	}

	return h, nil
}
