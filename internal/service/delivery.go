package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CamilleOnoda/webhook-relay.git/internal/database"
)

type DeliveryResult struct {
	StatusCode         sql.NullInt32
	ResponseBody       sql.NullString
	ErrorMessage       sql.NullString
	DeliveryDurationMs sql.NullInt32
}

func AttemptDelivery(ctx context.Context, event database.WebhookEvent, targetURL string) (DeliveryResult, error) {
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(outgoingReq)
	if err != nil {
		return DeliveryResult{
			ErrorMessage: sql.NullString{
				String: err.Error(),
				Valid:  true,
			},
		}, nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return DeliveryResult{}, fmt.Errorf("failed to read response body", err)
	}

	respBody := sql.NullString{
		String: string(bodyBytes),
		Valid:  true,
	}
	statusCode := sql.NullInt32{
		Int32: int32(resp.StatusCode),
		Valid: true,
	}
	result := DeliveryResult{
		StatusCode:   statusCode,
		ResponseBody: respBody,
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
