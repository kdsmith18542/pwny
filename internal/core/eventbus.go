// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"sync"
	"time"
)

type EventType string

const (
	EvtModuleStarted   EventType = "module:started"
	EvtModuleCompleted EventType = "module:completed"
	EvtModuleFailed    EventType = "module:failed"
	EvtSessionOpened   EventType = "session:opened"
	EvtSessionClosed   EventType = "session:closed"
	EvtJobCreated      EventType = "job:created"
	EvtJobUpdated      EventType = "job:updated"
)

type Event struct {
	Type EventType              `json:"type"`
	Time time.Time              `json:"time"`
	Data map[string]interface{} `json:"data"`
}

type EventBus struct {
	mu       sync.RWMutex
	channels []chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (eb *EventBus) Subscribe(buffer int) chan Event {
	ch := make(chan Event, buffer)
	eb.mu.Lock()
	eb.channels = append(eb.channels, ch)
	eb.mu.Unlock()
	return ch
}

func (eb *EventBus) Unsubscribe(ch chan Event) {
	eb.mu.Lock()
	for i, c := range eb.channels {
		if c == ch {
			eb.channels = append(eb.channels[:i], eb.channels[i+1:]...)
			close(ch)
			break
		}
	}
	eb.mu.Unlock()
}

func (eb *EventBus) Publish(evt Event) {
	eb.mu.RLock()
	channels := append([]chan Event{}, eb.channels...)
	eb.mu.RUnlock()
	for _, ch := range channels {
		select {
		case ch <- evt:
		default:
		}
	}
}
