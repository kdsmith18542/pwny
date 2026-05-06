// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kdsmith18542/pwny/internal/core"
	"github.com/kdsmith18542/pwny/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	core.Register("test_runner", func() core.Module {
		m := core.NewBaseModule(core.TypeAuxiliary, "test_runner")
		m.SetDescription("A test module for API tests")
		m.RegisterOption("MSG", "A message option", false, "hello")
		m.RegisterOption("REQUIRED_OPT", "A required option", true, "")
		return m
	})
}

func newTestServer(t *testing.T) *Server {
	t.Helper()
	eventBus := core.NewEventBus()
	sm := core.NewSessionManager()
	jm := core.NewJobManager(eventBus)
	s := New(Config{Host: "127.0.0.1", Port: 0, Allowed: []string{}}, sm, jm, eventBus, nil)
	return s
}

func TestStatusEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestListModulesEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/modules", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestListModulesEndpointWithTypeFilter(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/modules?type=auxiliary", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestGetModuleEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/modules/test_runner", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)

	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "test_runner", data["name"])
	assert.Equal(t, "auxiliary", data["type"])
}

func TestGetModuleEndpointNotFound(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/modules/nonexistent", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
}

func TestValidateModuleEndpoint(t *testing.T) {
	s := newTestServer(t)
	body := map[string]interface{}{
		"options": map[string]interface{}{
			"MSG":          "test message",
			"REQUIRED_OPT": "present",
		},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/modules/test_runner/validate", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestRunModuleEndpoint(t *testing.T) {
	s := newTestServer(t)
	body := map[string]interface{}{
		"options": map[string]interface{}{
			"REQUIRED_OPT": "present",
		},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/modules/test_runner/run", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "accepted", resp.Status)

	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, dataMap["job_id"])
}

func TestRunModuleEndpointNonexistent(t *testing.T) {
	s := newTestServer(t)
	body := map[string]interface{}{
		"options": map[string]interface{}{},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/modules/nonexistent/run", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListSessionsEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/sessions", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestGetSessionEndpointNotFound(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/sessions/nonexistent", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
}

func TestCloseSessionEndpointNotFound(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("DELETE", "/api/v1/sessions/nonexistent", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListJobsEndpoint(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/jobs", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestGetJobEndpointNotFound(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("GET", "/api/v1/jobs/nonexistent", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCancelJobEndpointNotFound(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/jobs/nonexistent/cancel", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJobCreatedByModuleRun(t *testing.T) {
	s := newTestServer(t)

	body := map[string]interface{}{
		"options": map[string]interface{}{
			"REQUIRED_OPT": "present",
		},
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/modules/test_runner/run", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	require.Equal(t, http.StatusAccepted, w.Code)

	var runResp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &runResp)
	require.NoError(t, err)
	dataMap := runResp.Data.(map[string]interface{})
	jobID := dataMap["job_id"].(string)

	getReq := httptest.NewRequest("GET", "/api/v1/jobs/"+jobID, nil)
	getW := httptest.NewRecorder()
	s.router.ServeHTTP(getW, getReq)

	assert.Equal(t, http.StatusOK, getW.Code)

	var jobResp APIResponse
	err = json.Unmarshal(getW.Body.Bytes(), &jobResp)
	require.NoError(t, err)
	assert.Equal(t, "ok", jobResp.Status)
}

func TestRecoveryMiddleware(t *testing.T) {
	s := newTestServer(t)

	s.router.Get("/api/v1/test-panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/api/v1/test-panic", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResp APIError
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", errResp.Error)
}

func TestCORSHeaders(t *testing.T) {
	s := New(Config{
		Host:    "127.0.0.1",
		Port:    0,
		Allowed: []string{"http://example.com"},
	}, core.NewSessionManager(), core.NewJobManager(core.NewEventBus()), core.NewEventBus(), nil)

	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSHeadersDenied(t *testing.T) {
	s := New(Config{
		Host:    "127.0.0.1",
		Port:    0,
		Allowed: []string{"http://example.com"},
	}, core.NewSessionManager(), core.NewJobManager(core.NewEventBus()), core.NewEventBus(), nil)

	req := httptest.NewRequest("GET", "/api/v1/status", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSPreflight(t *testing.T) {
	s := New(Config{
		Host:    "127.0.0.1",
		Port:    0,
		Allowed: []string{"http://example.com"},
	}, core.NewSessionManager(), core.NewJobManager(core.NewEventBus()), core.NewEventBus(), nil)

	req := httptest.NewRequest("OPTIONS", "/api/v1/status", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestValidateModuleNonexistent(t *testing.T) {
	s := newTestServer(t)
	body := map[string]interface{}{"options": map[string]interface{}{}}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/modules/nonexistent/validate", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestModuleValidationBadBody(t *testing.T) {
	s := newTestServer(t)
	req := httptest.NewRequest("POST", "/api/v1/modules/test_runner/validate", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSessionManagerInServer(t *testing.T) {
	s := newTestServer(t)
	assert.NotNil(t, s.sessions)
}

func TestJobManagerInServer(t *testing.T) {
	s := newTestServer(t)
	assert.NotNil(t, s.jobs)
}

func TestEventBusInServer(t *testing.T) {
	s := newTestServer(t)
	assert.NotNil(t, s.events)
}

func newTestServerWithDB(t *testing.T) *Server {
	t.Helper()
	dbPath := os.TempDir() + "/pwny_api_test_" + t.Name() + ".db"
	database, err := db.Open(dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		database.Close()
		os.Remove(dbPath)
	})

	eventBus := core.NewEventBus()
	sm := core.NewSessionManager()
	jm := core.NewJobManager(eventBus)
	s := New(Config{Host: "127.0.0.1", Port: 0, Allowed: []string{}}, sm, jm, eventBus, database)
	return s
}

func TestListWorkspacesEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	req := httptest.NewRequest("GET", "/api/v1/workspaces", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestCreateWorkspaceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	body := map[string]interface{}{"name": "test-workspace", "description": "test"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/workspaces", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestCreateWorkspaceEndpointNoName(t *testing.T) {
	s := newTestServerWithDB(t)
	body := map[string]interface{}{"description": "no name"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/workspaces", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetWorkspaceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)

	created := createTestWorkspace(t, s, "get-test-ws")
	req := httptest.NewRequest("GET", "/api/v1/workspaces/"+created, nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestGetWorkspaceEndpointNotFound(t *testing.T) {
	s := newTestServerWithDB(t)
	req := httptest.NewRequest("GET", "/api/v1/workspaces/nonexistent", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteWorkspaceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	created := createTestWorkspace(t, s, "delete-test-ws")

	req := httptest.NewRequest("DELETE", "/api/v1/workspaces/"+created, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListHostsInWorkspaceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	wsID := createTestWorkspace(t, s, "host-list-ws")

	req := httptest.NewRequest("GET", "/api/v1/workspaces/"+wsID+"/hosts", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
}

func TestCreateHostInWorkspaceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	wsID := createTestWorkspace(t, s, "host-create-ws")

	body := map[string]interface{}{"address": "10.0.0.1"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/workspaces/"+wsID+"/hosts", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateHostInMissingWorkspace(t *testing.T) {
	s := newTestServerWithDB(t)
	body := map[string]interface{}{"address": "10.0.0.1"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/workspaces/nonexistent/hosts", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHostEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	hostID := createTestHost(t, s, "get-host-ws", "10.0.0.5")

	req := httptest.NewRequest("GET", "/api/v1/hosts/"+hostID, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateHostEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	hostID := createTestHost(t, s, "update-host-ws", "10.0.0.6")

	body := map[string]interface{}{"os_name": "Windows", "state": "dead"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/api/v1/hosts/"+hostID, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteHostEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	hostID := createTestHost(t, s, "delete-host-ws", "10.0.0.7")

	req := httptest.NewRequest("DELETE", "/api/v1/hosts/"+hostID, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListServicesForHostEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	hostID := createTestHost(t, s, "svc-list-ws", "10.0.0.8")

	req := httptest.NewRequest("GET", "/api/v1/hosts/"+hostID+"/services", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateServiceForHostEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	hostID := createTestHost(t, s, "svc-create-ws", "10.0.0.9")

	body := map[string]interface{}{"port": 80, "proto": "tcp"}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/hosts/"+hostID+"/services", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetServiceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	serviceID := createTestService(t, s, "get-svc-ws", "10.0.0.10", 443)

	req := httptest.NewRequest("GET", "/api/v1/services/"+serviceID, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteServiceEndpoint(t *testing.T) {
	s := newTestServerWithDB(t)
	serviceID := createTestService(t, s, "delete-svc-ws", "10.0.0.11", 8080)

	req := httptest.NewRequest("DELETE", "/api/v1/services/"+serviceID, nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func createTestWorkspace(t *testing.T, s *Server, name string) string {
	t.Helper()
	w, err := s.database.CreateWorkspace(name, "test")
	require.NoError(t, err)
	return w.ID
}

func createTestHost(t *testing.T, s *Server, wsName, addr string) string {
	t.Helper()
	wsID := createTestWorkspace(t, s, wsName)
	h, err := s.database.CreateHost(wsID, addr)
	require.NoError(t, err)
	return h.ID
}

func createTestService(t *testing.T, s *Server, wsName, addr string, port int) string {
	t.Helper()
	hostID := createTestHost(t, s, wsName, addr)
	svc, err := s.database.CreateService(hostID, port, "tcp")
	require.NoError(t, err)
	return svc.ID
}
