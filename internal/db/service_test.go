// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateService(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")

	s, err := d.CreateService(h.ID, 80, "tcp")
	require.NoError(t, err)
	assert.NotEmpty(t, s.ID)
	assert.Equal(t, h.ID, s.HostID)
	assert.Equal(t, 80, s.Port)
	assert.Equal(t, "tcp", s.Proto)
	assert.Equal(t, "open", s.State)
}

func TestGetService(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")
	created, _ := d.CreateService(h.ID, 443, "tcp")

	got, err := d.GetService(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, 443, got.Port)
	assert.Equal(t, "tcp", got.Proto)
}

func TestGetServiceNotFound(t *testing.T) {
	d := newTestDB(t)
	_, err := d.GetService("nonexistent")
	assert.Error(t, err)
}

func TestListServices(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")

	d.CreateService(h.ID, 22, "tcp")
	d.CreateService(h.ID, 80, "tcp")
	d.CreateService(h.ID, 443, "tcp")

	services, err := d.ListServices(h.ID)
	require.NoError(t, err)
	assert.Len(t, services, 3)
}

func TestListServicesEmpty(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")

	services, err := d.ListServices(h.ID)
	require.NoError(t, err)
	assert.Empty(t, services)
}

func TestDeleteService(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")
	s, _ := d.CreateService(h.ID, 8080, "tcp")

	err := d.DeleteService(s.ID)
	require.NoError(t, err)

	_, err = d.GetService(s.ID)
	assert.Error(t, err)
}

func TestUpdateService(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")
	s, _ := d.CreateService(h.ID, 80, "tcp")

	updated, err := d.UpdateService(s.ID, func(svc *Service) {
		svc.Name = "http"
		svc.State = "filtered"
	})
	require.NoError(t, err)
	assert.Equal(t, "http", updated.Name)
	assert.Equal(t, "filtered", updated.State)

	got, _ := d.GetService(s.ID)
	assert.Equal(t, "http", got.Name)
}

func TestUpdateServiceNotFound(t *testing.T) {
	d := newTestDB(t)
	_, err := d.UpdateService("nonexistent", func(s *Service) {})
	assert.Error(t, err)
}
