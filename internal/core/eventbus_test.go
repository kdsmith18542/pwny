// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventBus(t *testing.T) {
	bus := NewEventBus()
	assert.NotNil(t, bus)
}

func TestEventBusSubscribePublish(t *testing.T) {
	bus := NewEventBus()
	ch := bus.Subscribe(10)

	bus.Publish(Event{
		Type: EvtModuleStarted,
		Time: time.Now(),
		Data: map[string]interface{}{"module": "test"},
	})

	select {
	case evt := <-ch:
		assert.Equal(t, EvtModuleStarted, evt.Type)
		assert.Equal(t, "test", evt.Data["module"])
	default:
		t.Fatal("expected event on channel")
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	bus := NewEventBus()
	ch := bus.Subscribe(10)
	bus.Unsubscribe(ch)

	bus.Publish(Event{
		Type: EvtModuleStarted,
		Time: time.Now(),
	})

	select {
	case _, ok := <-ch:
		assert.False(t, ok, "channel should be closed after unsubscribe")
	default:
	}

	bus.Publish(Event{
		Type: EvtModuleCompleted,
		Time: time.Now(),
	})
}

func TestEventBusMultipleSubscribers(t *testing.T) {
	bus := NewEventBus()
	ch1 := bus.Subscribe(10)
	ch2 := bus.Subscribe(10)

	bus.Publish(Event{
		Type: EvtJobCreated,
		Time: time.Now(),
		Data: map[string]interface{}{"job_id": "123"},
	})

	select {
	case evt := <-ch1:
		assert.Equal(t, "123", evt.Data["job_id"])
	default:
		t.Fatal("ch1 expected event")
	}

	select {
	case evt := <-ch2:
		assert.Equal(t, "123", evt.Data["job_id"])
	default:
		t.Fatal("ch2 expected event")
	}
}

func TestEventBusBufferFullDrops(t *testing.T) {
	bus := NewEventBus()
	ch := bus.Subscribe(1)

	bus.Publish(Event{Type: EvtModuleStarted, Time: time.Now()})
	bus.Publish(Event{Type: EvtModuleCompleted, Time: time.Now()})
	bus.Publish(Event{Type: EvtModuleFailed, Time: time.Now()})

	count := 0
	for {
		select {
		case <-ch:
			count++
		default:
			goto done
		}
	}
done:
	assert.Equal(t, 1, count)
}

func TestEventBusPublishNoSubscribers(t *testing.T) {
	bus := NewEventBus()
	bus.Publish(Event{Type: EvtSessionOpened, Time: time.Now()})
}
