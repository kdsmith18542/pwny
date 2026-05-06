// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJobManager(t *testing.T) {
	bus := NewEventBus()
	jm := NewJobManager(bus)
	assert.NotNil(t, jm)
	assert.Empty(t, jm.List())
}

func TestJobCreate(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", map[string]interface{}{"opt1": "val1"})

	assert.NotEmpty(t, job.ID)
	assert.Equal(t, "test/module", job.ModuleName)
	assert.Equal(t, JobPending, job.Status)
	assert.Equal(t, "val1", job.Options["opt1"])
	assert.False(t, job.CreatedAt.IsZero())
	assert.Nil(t, job.StartedAt)
	assert.Nil(t, job.CompletedAt)
}

func TestJobGet(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)

	got, err := jm.Get(job.ID)
	require.NoError(t, err)
	assert.Equal(t, job.ID, got.ID)
	assert.Equal(t, "test/module", got.ModuleName)
}

func TestJobGetNotFound(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	_, err := jm.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestJobStart(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)

	err := jm.Start(job.ID)
	require.NoError(t, err)

	got, _ := jm.Get(job.ID)
	assert.Equal(t, JobRunning, got.Status)
	assert.NotNil(t, got.StartedAt)
}

func TestJobStartNotFound(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	err := jm.Start("nonexistent")
	assert.Error(t, err)
}

func TestJobComplete(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)
	jm.Start(job.ID)
	jm.Complete(job.ID, "result data")

	got, _ := jm.Get(job.ID)
	assert.Equal(t, JobCompleted, got.Status)
	assert.Equal(t, "result data", got.Result)
	assert.NotNil(t, got.CompletedAt)
}

func TestJobFail(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)
	jm.Start(job.ID)
	jm.Fail(job.ID, assert.AnError)

	got, _ := jm.Get(job.ID)
	assert.Equal(t, JobFailed, got.Status)
	assert.Contains(t, got.Error, assert.AnError.Error())
	assert.NotNil(t, got.CompletedAt)
}

func TestJobCancel(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)

	err := jm.Cancel(job.ID)
	require.NoError(t, err)

	got, _ := jm.Get(job.ID)
	assert.Equal(t, JobCancelled, got.Status)
	assert.NotNil(t, got.CompletedAt)
}

func TestJobCancelCompleted(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	job := jm.Create("test/module", nil)
	jm.Start(job.ID)
	jm.Complete(job.ID, "ok")

	err := jm.Cancel(job.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not cancellable")
}

func TestJobList(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	j1 := jm.Create("mod1", nil)
	j2 := jm.Create("mod2", nil)

	jobs := jm.List()
	assert.Len(t, jobs, 2)

	ids := map[string]bool{j1.ID: true, j2.ID: true}
	for _, j := range jobs {
		delete(ids, j.ID)
	}
	assert.Empty(t, ids)
}

func TestJobEventsPublished(t *testing.T) {
	bus := NewEventBus()
	events := bus.Subscribe(10)

	jm := NewJobManager(bus)

	job := jm.Create("test/module", nil)
	evt := <-events
	assert.Equal(t, EvtJobCreated, evt.Type)
	assert.Equal(t, job.ID, evt.Data["job_id"])

	jm.Start(job.ID)
	evt = <-events
	assert.Equal(t, EvtModuleStarted, evt.Type)

	jm.Complete(job.ID, "ok")
	evt = <-events
	assert.Equal(t, EvtModuleCompleted, evt.Type)
}

func TestJobManagerWithoutBus(t *testing.T) {
	jm := NewJobManager(nil)
	job := jm.Create("test", nil)
	assert.NotNil(t, job)
	assert.Equal(t, JobPending, job.Status)
}

func TestJobIDUnique(t *testing.T) {
	jm := NewJobManager(NewEventBus())
	j1 := jm.Create("mod1", nil)
	j2 := jm.Create("mod2", nil)
	assert.NotEqual(t, j1.ID, j2.ID)
}

func TestGenerateJobID(t *testing.T) {
	id1 := generateJobID()
	id2 := generateJobID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 16)
}
