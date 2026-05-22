# Webhook Relay Service

A webhook relay service written in Go. It allows users to create webhook endpoints, 
receive incoming webhook events,<br>store them in PostgreSQL, 
and forward them to a configured target URL.

This project demonstrates backend fundamentals such as HTTP routing,
 JSON APIs, database persistence, webhook handling,<br>request forwarding, 
 delivery tracking, and basic frontend integration.

It focuses heavily on HTTP fundamentals and webhook forwarding behavior.

<table>
  <tr>
    <td align="center">
      <img src="assets/dashboard.png" height="350">
    </td>
    <td align="center">
      <img src="assets/event_delivery.png" height="350">
    </td>
  </tr>
</table>

## Webhook flow and Demo
1. A client registers a webhook endpoint with:
   `POST /api/endpoints`
2. The relay generates a unique webhook URL:
   `/webhooks/{endpoint_id}`
3. External services send webhook events to that generated URL
4. The relay validates the endpoint ID
5. The incoming payload and headers are stored in PostgreSQL
6. Unsafe or hop-by-hop headers are filtered out
7. The webhook request is forwarded to the configured target URL
8. Delivery results (status, duration, errors, response data) are stored for later inspection

<table>
  <tr>
    <td align="center">
      <img src="assets/webhook-diagram.png" height="350" width="350">
    </td>
    <td align="center">
      <img src="assets/webhook-relay-demo.gif" height="350" width="350">
    </td>
  </tr>
</table>

## Prerequisites and setup

This project requires:

- Go
- PostgreSQL
- Goose for database migrations
- sqlc for generated database code

### Clone the repository

```bash
git clone https://github.com/CamilleOnoda/webhook-relay.git
cd webhook-relay
```

### Install dependencies
```bash
go mod download
```

### Set up environment variables
Create a .env file:
```bash
DB_URL=postgres://username:password@localhost:5432/webhook_relay?sslmode=disable
PORT=8080
BASE_URL=http://localhost:8080
```

### Run database migrations
```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### Generate sqlc code
```bash
sqlc generate
```

### Run the server
```bash
go run .
```
The API should be available at:
```bash
http://localhost:8080
```

## Project overview
This service lets you:

- Create webhook endpoints
- Generate a public webhook URL for each endpoint
- Receive webhook requests
- Store webhook payloads and headers
- Forward webhook events to a target URL
- Track delivery status and response data
- View endpoints and delivery history
- Delete endpoints and related webhook data

## API features
#### Health check
```bash
GET /api/health
```
Returns a simple readiness response.

#### Create a webhook endpoint
```bash
curl -X POST http://localhost:8080/api/endpoints \
-H "Content-Type: application/json" \
-d '{
  "name": "Test endpoint",
  "target_url": "https://httpbin.org/post"
}'
```
Example response:

```json
{
  "id": "6d8e487d-1cd8-4bed-96a5-fbc98cde79be",
  "name": "httpbin success test",
  "target_url": "https://httpbin.org/post",
  "generated_url": "http://localhost:8080/webhooks/6d8e487d-1cd8-4bed-96a5-fbc98cde79be"
}
```

#### List endpoints
```bash
curl http://localhost:8080/api/endpoints
```
Returns all configured webhook endpoints.

#### Get endpoint by ID
```bash
GET /api/endpoints/{id}
```
Returns a single endpoint.

#### Delete endpoint by ID
```bash
DELETE /api/endpoints/{id}
```
Deletes an endpoint and its related webhook data.

#### Delete all endpoints
```bash
curl -X DELETE http://localhost:8080/admin/endpoints/delete
```
Deletes all endpoints and related webhook data.

#### Receive webhook
```bash
curl -X POST http://localhost:8080/webhooks/6d8e487d-1cd8-4bed-96a5-fbc98cde79be \
-H "Content-Type: application/json" \
-H "X-Event-Type: user.created" \
-d '{
  "message": "hello"
}'
```
Receives an incoming webhook event, stores the payload and headers, then attempts delivery to the endpoint target URL.

#### List deliveries
```bash
curl http://localhost:8080/api/deliveries
```
Returns webhook delivery history.

#### Example failing delivery

Create an endpoint pointing to an HTTP 500 target:

```bash
curl -X POST http://localhost:8080/api/endpoints \
-H "Content-Type: application/json" \
-d '{
  "name": "httpbin failure test",
  "target_url": "https://httpbin.org/status/500",
  "description": "500 status code test"
}'
```

## Frontend

A lightweight frontend is included for endpoint management and manual webhook testing.
It uses:

- HTML
- CSS
- Vanilla JavaScript
- fetch()

#### The frontend can:

- List endpoints
- Create endpoints
- Delete endpoints
- Send test webhooks
- View delivery/event history

## Architecture and design
- main.go: Starts the server and registers routes
- json.go: Shared JSON response helpers
- readiness.go: Health check handler
- handler_endpoint_create.go: Creates webhook endpoints
- handler_endpoints_get.go: Lists endpoints
- handler_endpointsByID_get.go: Gets one endpoint by ID
- handler_endpointByID_delete.go: Deletes one endpoint
- handler_webhook_receive.go: Receives incoming webhook requests
- handler_deliveries_get.go: Lists webhook delivery records
- internal/database: sqlc-generated database access code
- internal/service: Webhook delivery and header filtering logic
- sql/schema: Database migrations
- static: Basic frontend files
- test_REST_client.http: Manual API testing requests

## Key patterns
- HTTP handlers separated by responsibility
- PostgreSQL persistence with sqlc-generated queries
- Database migrations with Goose
- UUID-based webhook endpoint URLs
- JSON request/response handling
- Request header filtering before forwarding
- Delivery status tracking
- Clear separation between API handlers, database access, and service logic
- Simple frontend for manual testing and demonstration

## What I learned
This project helped me better understand:
- HTTP request/response lifecycle
- Webhook delivery patterns
- JSON payload handling
- Forwarding requests safely
- Hop-by-hop and unsafe headers
- PostgreSQL persistence with sqlc
- Database migrations with Goose
- UUID-based routing
- Delivery tracking and retry-oriented design

## Post MVP roadmap
Improvements after the initial MVP:

- Authentication and user accounts
  - Associate endpoints with users
  - Protect endpoint management routes

- Webhook signature verification
  - Verify incoming requests using shared secrets
  - Reject spoofed or unauthorized webhook events

- Queue-based delivery system
  - Move delivery attempts to a background worker
  - Avoid blocking the webhook receiver while forwarding requests

- Retry logic
  - Retry failed deliveries with backoff
  - Track retry attempts and final delivery state

- Concurrency control
  - Limit simultaneous delivery attempts
  - Prevent one slow target from affecting the whole service

- Endpoint secrets
  - Generate a secret per endpoint
  - Use it for signing or validating webhook requests

- Delivery filtering
  - Allow endpoints to subscribe only to specific event types

- Better observability
  - Add structured logs
  - Add metrics for successful and failed deliveries