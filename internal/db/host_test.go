// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *Database {
	t.Helper()
	db, err := Open(os.TempDir() + "/pwny_test_" + t.Name() + ".db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Close()
		os.Remove(os.TempDir() + "/pwny_test_" + t.Name() + ".db")
	})
	return db
}

func TestCreateHost(t *testing.T) {
	d := newTestDB(t)
	w, err := d.CreateWorkspace("test", "test workspace")
	require.NoError(t, err)

	h, err := d.CreateHost(w.ID, "192.168.1.1")
	require.NoError(t, err)
	assert.NotEmpty(t, h.ID)
	assert.Equal(t, w.ID, h.WorkspaceID)
	assert.Equal(t, "192.168.1.1", h.Address)
	assert.Equal(t, "alive", h.State)
	assert.False(t, h.CreatedAt.IsZero())
}

func TestGetHost(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	created, _ := d.CreateHost(w.ID, "10.0.0.1")

	got, err := d.GetHost(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "10.0.0.1", got.Address)
}

func TestGetHostNotFound(t *testing.T) {
	d := newTestDB(t)
	_, err := d.GetHost("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestListHosts(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")

	d.CreateHost(w.ID, "10.0.0.1")
	d.CreateHost(w.ID, "10.0.0.2")

	hosts, err := d.ListHosts(w.ID)
	require.NoError(t, err)
	assert.Len(t, hosts, 2)
}

func TestListHostsEmpty(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")

	hosts, err := d.ListHosts(w.ID)
	require.NoError(t, err)
	assert.Empty(t, hosts)
}

func TestDeleteHost(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")

	err := d.DeleteHost(h.ID)
	require.NoError(t, err)

	_, err = d.GetHost(h.ID)
	assert.Error(t, err)
}

func TestUpdateHost(t *testing.T) {
	d := newTestDB(t)
	w, _ := d.CreateWorkspace("test", "")
	h, _ := d.CreateHost(w.ID, "10.0.0.1")

	updated, err := d.UpdateHost(h.ID, func(host *Host) {
		host.OSName = "Linux"
		host.State = "dead"
	})
	require.NoError(t, err)
	assert.Equal(t, "Linux", updated.OSName)
	assert.Equal(t, "dead", updated.State)

	got, _ := d.GetHost(h.ID)
	assert.Equal(t, "Linux", got.OSName)
}

func TestUpdateHostNotFound(t *testing.T) {
	d := newTestDB(t)
	_, err := d.UpdateHost("nonexistent", func(h *Host) {})
	assert.Error(t, err)
}
