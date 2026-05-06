// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ID        string    `json:"id"`
	HostID    string    `json:"host_id"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Name      string    `json:"name,omitempty"`
	Info      string    `json:"info,omitempty"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const svcCols = `id, host_id, port, proto, COALESCE(name,''), COALESCE(info,''), state, created_at, updated_at`

func (d *Database) CreateService(hostID string, port int, proto string) (*Service, error) {
	s := &Service{
		ID:        uuid.New().String(),
		HostID:    hostID,
		Port:      port,
		Proto:     proto,
		State:     "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := d.Exec(
		`INSERT INTO services (id, host_id, port, proto, state, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		s.ID, s.HostID, s.Port, s.Proto, s.State, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return s, nil
}

func (d *Database) GetService(id string) (*Service, error) {
	s := &Service{}
	err := d.QueryRow(
		`SELECT `+svcCols+` FROM services WHERE id = ?`, id,
	).Scan(&s.ID, &s.HostID, &s.Port, &s.Proto, &s.Name, &s.Info, &s.State, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("service not found: %s", id)
	}
	return s, nil
}

func (d *Database) ListServices(hostID string) ([]Service, error) {
	rows, err := d.Query(
		`SELECT `+svcCols+` FROM services WHERE host_id = ? ORDER BY port`, hostID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		if err := rows.Scan(&s.ID, &s.HostID, &s.Port, &s.Proto, &s.Name, &s.Info, &s.State, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, rows.Err()
}

func (d *Database) DeleteService(id string) error {
	_, err := d.Exec(`DELETE FROM services WHERE id = ?`, id)
	return err
}

func (d *Database) UpdateService(id string, updater func(s *Service)) (*Service, error) {
	s, err := d.GetService(id)
	if err != nil {
		return nil, err
	}

	updater(s)
	s.UpdatedAt = time.Now()

	_, err = d.Exec(
		`UPDATE services SET name=?, info=?, state=?, updated_at=? WHERE id=?`,
		s.Name, s.Info, s.State, s.UpdatedAt, s.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return s, nil
}
