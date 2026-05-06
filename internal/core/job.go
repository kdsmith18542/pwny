// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
	JobCancelled JobStatus = "cancelled"
)

type Job struct {
	ID          string                 `json:"id"`
	ModuleName  string                 `json:"module_name"`
	Status      JobStatus              `json:"status"`
	Options     map[string]interface{} `json:"options"`
	Result      interface{}            `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

type JobManager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
	bus  *EventBus
}

func NewJobManager(bus *EventBus) *JobManager {
	return &JobManager{
		jobs: make(map[string]*Job),
		bus:  bus,
	}
}

func (jm *JobManager) Create(moduleName string, options map[string]interface{}) *Job {
	job := &Job{
		ID:         generateJobID(),
		ModuleName: moduleName,
		Status:     JobPending,
		Options:    options,
		CreatedAt:  time.Now(),
	}

	jm.mu.Lock()
	jm.jobs[job.ID] = job
	jm.mu.Unlock()

	slog.Info("job created", "job_id", job.ID, "module", moduleName)

	if jm.bus != nil {
		jm.bus.Publish(Event{
			Type: EvtJobCreated,
			Time: time.Now(),
			Data: map[string]interface{}{
				"job_id":  job.ID,
				"module":  moduleName,
				"status":  string(job.Status),
				"options": options,
			},
		})
	}

	return job
}

func (jm *JobManager) Get(id string) (*Job, error) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	job, exists := jm.jobs[id]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", id)
	}

	return job, nil
}

func (jm *JobManager) List() []*Job {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	jobs := make([]*Job, 0, len(jm.jobs))
	for _, j := range jm.jobs {
		jobs = append(jobs, j)
	}

	return jobs
}

func (jm *JobManager) Start(id string) error {
	jm.mu.Lock()
	job, exists := jm.jobs[id]
	if !exists {
		jm.mu.Unlock()
		return fmt.Errorf("job not found: %s", id)
	}

	now := time.Now()
	job.Status = JobRunning
	job.StartedAt = &now
	jm.mu.Unlock()

	slog.Info("job started", "job_id", id, "module", job.ModuleName)
	jm.publishUpdate(job)
	return nil
}

func (jm *JobManager) Complete(id string, result interface{}) {
	jm.mu.Lock()
	job, exists := jm.jobs[id]
	if !exists {
		jm.mu.Unlock()
		return
	}

	now := time.Now()
	job.Status = JobCompleted
	job.Result = result
	job.CompletedAt = &now
	jm.mu.Unlock()

	slog.Info("job completed", "job_id", id, "module", job.ModuleName)
	jm.publishUpdate(job)
}

func (jm *JobManager) Fail(id string, err error) {
	jm.mu.Lock()
	job, exists := jm.jobs[id]
	if !exists {
		jm.mu.Unlock()
		return
	}

	now := time.Now()
	job.Status = JobFailed
	job.Error = err.Error()
	job.CompletedAt = &now
	jm.mu.Unlock()

	slog.Info("job failed", "job_id", id, "module", job.ModuleName, "error", err)
	jm.publishUpdate(job)
}

func (jm *JobManager) Cancel(id string) error {
	jm.mu.Lock()
	job, exists := jm.jobs[id]
	if !exists {
		jm.mu.Unlock()
		return fmt.Errorf("job not found: %s", id)
	}

	if job.Status != JobPending && job.Status != JobRunning {
		jm.mu.Unlock()
		return fmt.Errorf("job %s is not cancellable (status: %s)", id, job.Status)
	}

	now := time.Now()
	job.Status = JobCancelled
	job.CompletedAt = &now
	jm.mu.Unlock()

	slog.Info("job cancelled", "job_id", id, "module", job.ModuleName)
	jm.publishUpdate(job)
	return nil
}

func (jm *JobManager) publishUpdate(job *Job) {
	if jm.bus == nil {
		return
	}

	eventType := EvtJobUpdated
	if job.Status == JobCompleted || job.Status == JobFailed || job.Status == JobCancelled {
		if job.Status == JobCompleted {
			eventType = EvtModuleCompleted
		} else if job.Status == JobFailed {
			eventType = EvtModuleFailed
		}
	} else if job.Status == JobRunning {
		eventType = EvtModuleStarted
	}

	jm.bus.Publish(Event{
		Type: eventType,
		Time: time.Now(),
		Data: map[string]interface{}{
			"job_id": job.ID,
			"module": job.ModuleName,
			"status": string(job.Status),
		},
	})
}

func generateJobID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		slog.Error("failed to generate job ID", "error", err)
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(buf)
}
