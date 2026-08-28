package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
	"github.com/google/uuid"
)

type DeliveryResult struct {
	StatusCode         sql.NullInt32
	ResponseBody       sql.NullString
	ErrorMessage       sql.NullString
	DeliveryDurationMs sql.NullInt32
}

type DBQuerier interface {
	GetReadyDeliveries(ctx context.Context, limit int32) ([]database.Delivery, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (database.WebhookEvent, error)
	UpdateDeliveryState(ctx context.Context, arg database.UpdateDeliveryStateParams) error
}

type DeliveryService struct {
	httpClient *http.Client
}

type DatabaseDeliveryProcessor struct {
	db              DBQuerier
	deliveryService *DeliveryService
	maxAttempts     int
}

func NewDeliveryService() *DeliveryService {
	return &DeliveryService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func NewDatabaseDeliveryProcessor(db DBQuerier, ds *DeliveryService, maxAttempts int) *DatabaseDeliveryProcessor {
	return &DatabaseDeliveryProcessor{
		db:              db,
		deliveryService: ds,
		maxAttempts:     maxAttempts,
	}
}

func (p *DatabaseDeliveryProcessor) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.ProcessPendingDeliveries(ctx)
		}
	}
}

func (p *DatabaseDeliveryProcessor) ProcessPendingDeliveries(ctx context.Context) {
	deliveries, err := p.db.GetReadyDeliveries(ctx, 20)
	if err != nil {
		log.Printf("error fetching ready deliveries: %v", err)
		return
	}

	for _, delivery := range deliveries {
		event, err := p.db.GetEventByID(ctx, delivery.EventID)
		if err != nil {
			log.Printf("failed to fetch event for delivery id=%s: %v", delivery.ID, err)
			continue
		}

		result, deliveryErr := p.deliveryService.AttemptDelivery(ctx, event, delivery.TargetUrl)

		statusCode := int(result.StatusCode.Int32)
		isRetryable := isRetryable(statusCode, deliveryErr)
		currentAttempt := delivery.AttemptCount + 1

		var deliveryStatus string
		var nextRetryAt sql.NullTime
		var deliveredAt sql.NullTime

		if deliveryErr == nil && statusCode >= 200 && statusCode < 300 {
			deliveryStatus = "success"
			deliveredAt = sql.NullTime{Time: time.Now(), Valid: true}
		} else if isRetryable && currentAttempt < int32(p.maxAttempts) {
			deliveryStatus = "retry_scheduled"
			nextDelay := p.calculateBackoff(currentAttempt)
			nextRetryAt = sql.NullTime{Time: time.Now().Add(nextDelay), Valid: true}
		} else {
			deliveryStatus = "dead_letter"
		}

		errorMessage := ""
		if deliveryErr != nil {
			errorMessage = deliveryErr.Error()
		} else if result.ErrorMessage.Valid {
			errorMessage = result.ErrorMessage.String
		}
		nextRetry := "none"
		if nextRetryAt.Valid {
			nextRetry = nextRetryAt.Time.Format(time.RFC3339)
		}
		err = p.db.UpdateDeliveryState(ctx, database.UpdateDeliveryStateParams{
			ID:                 delivery.ID,
			Status:             deliveryStatus,
			StatusCode:         result.StatusCode,
			ResponseBody:       result.ResponseBody,
			ErrorMessage:       result.ErrorMessage,
			NextRetryAt:        nextRetryAt,
			DeliveredAt:        deliveredAt,
			DeliveryDurationMs: result.DeliveryDurationMs,
		})
		if err != nil {
			log.Printf("delivery attempt state update failed: delivery_id=%s | endpoint=%q | attempt=%d | intended_status=%s | http=%d | error=%q | retry=%s | persisted=false | db_error=%v",
				delivery.ID, delivery.TargetUrl, currentAttempt, deliveryStatus, statusCode, errorMessage, nextRetry, err)
			continue
		}

		log.Printf("delivery attempt completed: delivery_id=%s | endpoint=%q | attempt=%d | status=%s | http=%d | error=%q | retry=%s | persisted=true",
			delivery.ID, delivery.TargetUrl, currentAttempt, deliveryStatus, statusCode, errorMessage, nextRetry)
	}
}

func (p *DatabaseDeliveryProcessor) calculateBackoff(attempt int32) time.Duration {
	initialDelay := 5 * time.Second
	maxDelay := 1 * time.Hour

	multiplier := int64(1) << uint(attempt-1)
	delay := initialDelay * time.Duration(multiplier)

	if delay > maxDelay {
		delay = maxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(delay / 10)))
	return delay + jitter
}

func isRetryable(statusCode int, err error) bool {
	if err != nil {
		return true
	}
	switch statusCode {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (s *DeliveryService) AttemptDelivery(ctx context.Context, event database.WebhookEvent, targetURL string) (DeliveryResult, error) {
	body := bytes.NewReader(event.Payload)
	outgoingReq, err := http.NewRequestWithContext(ctx, "POST", targetURL, body)
	if err != nil {
		return DeliveryResult{}, fmt.Errorf("failed to create outgoing request: %w", err)
	}

	var originalHeaders map[string][]string
	if err := json.Unmarshal(event.Headers, &originalHeaders); err != nil {
		return DeliveryResult{}, fmt.Errorf("failed to unmarshal headers: %w", err)
	}

	for key, value := range originalHeaders {
		canonicalKey := http.CanonicalHeaderKey(key)
		if isUnsafeHeader(canonicalKey) {
			continue
		}
		for _, v := range value {
			outgoingReq.Header.Add(canonicalKey, v)
		}
	}

	var eventType string
	if event.EventType.Valid {
		eventType = event.EventType.String
	} else {
		eventType = "unknown"
	}

	outgoingReq.Header.Set("Content-Type", "application/json")
	outgoingReq.Header.Set("X-Relay-ID", event.ID.String())
	outgoingReq.Header.Set("X-Relay-Event-Type", eventType)
	outgoingReq.Header.Set("X-Relay-Received-At", event.ReceivedAt.Format(http.TimeFormat))

	start := time.Now()

	resp, err := s.httpClient.Do(outgoingReq)
	if err != nil {
		return DeliveryResult{
			ErrorMessage: sql.NullString{String: err.Error(), Valid: true},
		}, err
	}
	defer resp.Body.Close()

	limitedReader := io.LimitReader(resp.Body, 4096)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return DeliveryResult{}, fmt.Errorf("failed to read response body: %w", err)
	}

	total := time.Since(start)

	respBody := sql.NullString{String: string(bodyBytes), Valid: true}
	statusCode := sql.NullInt32{Int32: int32(resp.StatusCode), Valid: true}
	duration := sql.NullInt32{Int32: int32(total.Milliseconds()), Valid: true}

	result := DeliveryResult{
		StatusCode:         statusCode,
		ResponseBody:       respBody,
		DeliveryDurationMs: duration,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		result.ErrorMessage = sql.NullString{
			String: fmt.Sprintf("received non-2xx response: %d", resp.StatusCode),
			Valid:  true,
		}
		return result, nil
	}

	return result, nil
}
