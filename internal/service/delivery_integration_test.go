package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/CamilleOnoda/webhook-relay.git/internal/service"
	"github.com/google/uuid"
)

type mockDB struct {
	mu            sync.Mutex
	deliveries    []database.Delivery
	events        map[uuid.UUID]database.WebhookEvent
	updatedParams []database.UpdateDeliveryStateParams
}

func (m *mockDB) GetReadyDeliveries(ctx context.Context, limit int32) ([]database.Delivery, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ready []database.Delivery
	now := time.Now()
	for _, d := range m.deliveries {
		if d.Status == "pending" || (d.Status == "retry_scheduled" && d.NextRetryAt.Valid && d.NextRetryAt.Time.Before(now)) {
			ready = append(ready, d)
			if int32(len(ready)) >= limit {
				break
			}
		}
	}
	return ready, nil
}

func (m *mockDB) GetEventByID(ctx context.Context, id uuid.UUID) (database.WebhookEvent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	event, exists := m.events[id]
	if !exists {
		return database.WebhookEvent{}, errors.New("event not found")
	}
	return event, nil
}

func (m *mockDB) UpdateDeliveryState(ctx context.Context, arg database.UpdateDeliveryStateParams) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.updatedParams = append(m.updatedParams, arg)

	for i, d := range m.deliveries {
		if d.ID == arg.ID {
			m.deliveries[i].Status = arg.Status
			m.deliveries[i].AttemptCount++
			m.deliveries[i].NextRetryAt = arg.NextRetryAt
			break
		}
	}
	return nil
}

func TestIntegration_DeliverySuccess(t *testing.T) {
	serverCallCount := 0
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCallCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"received"}`))
	}))
	defer targetServer.Close()

	eventID := uuid.New()
	deliveryID := uuid.New()

	mock := &mockDB{
		events: map[uuid.UUID]database.WebhookEvent{
			eventID: {
				ID:        eventID,
				Payload:   []byte(`{"user":"test"}`),
				Headers:   []byte(`{"Content-Type":["application/json"]}`),
				EventType: sql.NullString{String: "user.created", Valid: true},
			},
		},
		deliveries: []database.Delivery{
			{
				ID:           deliveryID,
				EventID:      eventID,
				TargetUrl:    targetServer.URL,
				Status:       "pending",
				AttemptCount: 0,
			},
		},
	}

	deliveryService := service.NewDeliveryService()
	processor := service.NewDatabaseDeliveryProcessor(mock, deliveryService, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	processor.ProcessPendingDeliveries(ctx)

	if serverCallCount != 1 {
		t.Errorf("Expected 1 call, got: %d", serverCallCount)
	}

	mock.mu.Lock()
	if len(mock.updatedParams) != 1 {
		t.Fatalf("database should have been updated once, got: %d", len(mock.updatedParams))
	}

	finalState := mock.updatedParams[0]
	mock.mu.Unlock()

	if finalState.Status != "success" {
		t.Errorf("Final status should be 'success', got: %s", finalState.Status)
	}
	if finalState.StatusCode.Int32 != 200 {
		t.Errorf("HTTP status code should be 200, got: %d", finalState.StatusCode.Int32)
	}
}

func TestIntegration_DeliveryRetryOn500(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer targetServer.Close()

	eventID := uuid.New()
	deliveryID := uuid.New()

	mock := &mockDB{
		events: map[uuid.UUID]database.WebhookEvent{
			eventID: {
				ID:      eventID,
				Payload: []byte(`{}`),
				Headers: []byte(`{}`),
			},
		},
		deliveries: []database.Delivery{
			{
				ID:           deliveryID,
				EventID:      eventID,
				TargetUrl:    targetServer.URL,
				Status:       "pending",
				AttemptCount: 0,
			},
		},
	}

	deliveryService := service.NewDeliveryService()
	processor := service.NewDatabaseDeliveryProcessor(mock, deliveryService, 3)

	ctx := context.Background()
	processor.ProcessPendingDeliveries(ctx)

	mock.mu.Lock()
	if len(mock.updatedParams) != 1 {
		t.Fatalf("database should have been updated once, got: %d", len(mock.updatedParams))
	}

	firstAttempt := mock.updatedParams[0]
	mock.mu.Unlock()

	if firstAttempt.Status != "retry_scheduled" {
		t.Errorf("Status should be 'retry_scheduled' after a 500 code, got: %s", firstAttempt.Status)
	}
	if !firstAttempt.NextRetryAt.Valid {
		t.Error("'next_retry_at' should have been calculated")
	}
	if firstAttempt.Status == "dead_letter" {
		t.Error("Status should not be dead_letter on the first failed attempt")
	}
}

func TestIntegration_DeliveryDeadLetterAfterMaxAttempt(t *testing.T) {
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer targetServer.Close()

	eventID := uuid.New()
	deliveryID := uuid.New()

	mock := &mockDB{
		events: map[uuid.UUID]database.WebhookEvent{
			eventID: {
				ID:      eventID,
				Payload: []byte(`{}`),
				Headers: []byte(`{}`),
			},
		},
		deliveries: []database.Delivery{
			{
				ID:           deliveryID,
				EventID:      eventID,
				TargetUrl:    targetServer.URL,
				Status:       "retry_scheduled",
				AttemptCount: 2,
				NextRetryAt:  sql.NullTime{Time: time.Now().Add(-1 * time.Minute), Valid: true},
			},
		},
	}

	deliveryService := service.NewDeliveryService()
	processor := service.NewDatabaseDeliveryProcessor(mock, deliveryService, 3)

	ctx := context.Background()
	processor.ProcessPendingDeliveries(ctx)

	mock.mu.Lock()
	if len(mock.updatedParams) != 1 {
		t.Fatalf("database should have been updated once, got: %d", len(mock.updatedParams))
	}

	finalAttempt := mock.updatedParams[0]
	mock.mu.Unlock()

	if finalAttempt.Status != "dead_letter" {
		t.Errorf("Expected status to be 'dead_letter' after max attempts reached, got: %s", finalAttempt.Status)
	}

	if finalAttempt.NextRetryAt.Valid {
		t.Error("A dead letter delivery should not have a valid 'next_retry_at' timestamp")
	}
}
