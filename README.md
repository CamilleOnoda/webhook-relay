# Webhook Relay Service

A webhook relay service written in Go.

This project allows users to create webhook endpoints, receive incoming webhook events, store them in PostgreSQL, <br>and forward them to configured target URLs while tracking delivery results and request metadata.

The project focuses heavily on backend fundamentals such as:

- HTTP routing and request handling
- Webhook forwarding behavior
- JSON APIs
- PostgreSQL persistence
- Authentication and authorization
- Header filtering and request safety
- Delivery tracking and service architecture

---

# Motivation

Many modern applications rely on webhooks to communicate between services, yet the mechanics behind webhook delivery are often hidden behind third-party platforms.

I built this project to better understand what happens between the moment a webhook is received and the moment it reaches its destination.

The goal was not only to learn Go, but to gain experience with:

- HTTP request handling
- API design
- Database persistence
- Authentication and authorization
- Request forwarding
- Delivery tracking
- Service architecture

Rather than consuming webhooks through existing tools, I wanted to build the infrastructure myself and explore the challenges involved in receiving, storing, processing, and forwarding webhook events reliably.

This project became an opportunity to combine several backend concepts into a single application while simulating the kind of event-driven systems commonly used in production environments.

---

## Table of Contents

- [Features](#features)
- [Webhook Flow](#webhook-flow)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Header Filtering](#header-filtering)
- [What I Learned](#what-i-learned)
- [Post-MVP Roadmap](#post-mvp-roadmap)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Screenshots

<table>
  <tr>
    <td align="center">
      <img src="assets/login.png" height="320">
    </td>
    <td align="center">
      <img src="assets/dashboard.png" height="350">
    </td>
  </tr>
</table>

---

# Features

## Webhook endpoints

- Create webhook endpoints
- Generate unique UUID-based webhook URLs
- Store endpoint metadata in PostgreSQL
- Delete endpoints and related webhook data

## Webhook processing

- Receive incoming webhook events
- Store payloads and headers
- Forward requests to external target URLs
- Filter unsafe and hop-by-hop headers
- Track delivery status and duration
- Store response metadata and errors

## Authentication

- User registration
- User login
- Password hashing
- JWT authentication
- Protected API routes
- Admin-only routes

## Frontend dashboard

A lightweight frontend is included for manual testing and endpoint management.

Built using:

- HTML
- CSS
- Vanilla JavaScript
- fetch()

The frontend supports:

- Creating endpoints
- Listing endpoints
- Deleting endpoints
- Sending test webhooks
- Viewing webhook events
- Viewing delivery history

---

# Webhook Flow

1. A user creates a webhook endpoint through the API
2. The relay generates a unique webhook URL:

```text
/webhooks/{endpoint_id}
```

3. External services send webhook events to that URL
4. The relay validates the endpoint ID
5. The incoming payload and headers are stored
6. Unsafe headers are filtered
7. The webhook is forwarded to the configured target URL
8. Delivery results are stored for later inspection

<table>
  <tr>
    <td align="center">
      <img src="assets/webhook-diagram.png" height="350" width="350">
    </td>
    <td align="center">
      <img src="assets/webhook.gif" height="350" width="350">
    </td>
  </tr>
</table>

---

# Tech Stack

## Backend

- Go
- net/http
- PostgreSQL
- sqlc
- Goose

## Frontend

- HTML
- CSS
- Vanilla JavaScript

## Infrastructure

- Railway
- PostgreSQL
- GitHub

---

# Project Structure

```text
├── assets
├── internal
│   ├── auth
│   ├── database
│   └── service
├── static
├── sql
│   └── schema
├── .env
├── .gitignore
├── admin_middleware.go
├── auth_middleware.go
├── create_endpoint_test.go
├── go.mod
├── handler_deliveries_get.go
├── handler_endpoint_create.go
├── handler_endpointByID_delete.go
├── handler_endpointByID_get.go
├── handler_endpoints_delete.go
├── handler_endpoints_get.go
├── handler_events_get.go
├── handler_login.go
├── handler_users_create.go
├── handler_users_delete.go
├── handler_webhook_receive.go
├── json.go
├── main.go
├── readiness.go
├── README.md
├── sqlc.yaml
├── test_REST_client.http
└── webhook-relay.git
```

## Important files
```text
| File | Responsibility |
|---|---|
| `main.go` | Starts the server and registers routes |
| `json.go` | Shared JSON response helpers |
| `readiness.go` | Health check handler |
| `auth_middleware.go` | JWT authentication middleware |
| `admin_middleware.go` | Admin authorization middleware |
| `handler_login.go` | User login and JWT generation |
| `handler_users_create.go` | User registration |
| `handler_users_delete.go` | Admin route for deleting users |
| `handler_endpoint_create.go` | Create a new webhook endpoint for the authenticated user |
| `handler_endpoints_get.go` | Lists webhook endpoints for the authenticated user |
| `handler_endpointByID_get.go` | Gets a specific endpoint by its ID for the authenticated user |
| `handler_endpointByID_delete.go` | Deletes a specific endpoint by ID for the authenticated user |
| `handler_endpoints_delete.go` | Admin endpoint cleanup route |
| `handler_webhook_receive.go` | Receives incoming webhook requests |
| `handler_events_get.go` | Lists all webhook events for the authenticated user |
| `handler_deliveries_get.go` | Lists delivery history for the authenticated user |
| `internal/auth` | Password hashing and JWT validation |
| `internal/database` | sqlc-generated database layer |
| `internal/service` | Webhook forwarding and header filtering |
| `sql/schema` | Goose database migrations |
| `static` | Frontend files |
| `test_REST_client.http` | Manual API testing scenarios |
```

---

# Quick Start
## Prerequisites

- Go
- PostgreSQL
- Goose
- sqlc


## Clone the repository

```bash
git clone https://github.com/CamilleOnoda/webhook-relay.git
cd webhook-relay
```

## Install dependencies

```bash
go mod download
```

---

## Environment variables

Create a `.env` file:

```env
DB_URL=postgres://username:password@localhost:5432/webhook_relay?sslmode=disable
PORT=8080
BASE_URL=http://localhost:8080
JWT_SECRET=your_secret
```

---

## Run database migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

---

## Generate sqlc code

```bash
sqlc generate
```

---

## Run the server

```bash
go run .
```

The API should now be available at:

```text
http://localhost:8080
```

---

# Usage

## Health check

```http
GET /api/health
```

Returns a readiness response.

---

## Create user

```http
POST /api/users
```

Example:

```json
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password"
}
```

---

## Login

```http
POST /api/login
```

Returns a JWT token.

---

## Create endpoint

```bash
curl -X POST http://localhost:8080/api/endpoints \
-H "Authorization: Bearer YOUR_TOKEN" \
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
  "name": "Test endpoint",
  "target_url": "https://httpbin.org/post",
  "generated_url": "http://localhost:8080/webhooks/6d8e487d-1cd8-4bed-96a5-fbc98cde79be"
}
```

---

## Receive webhook

```bash
curl -X POST http://localhost:8080/webhooks/{endpoint_id} \
-H "Content-Type: application/json" \
-H "X-Event-Type: user.created" \
-d '{
  "message": "hello"
}'
```

The relay stores the event and forwards it to the configured target URL.

---

## List endpoints

```http
GET /api/endpoints
```

---

## List deliveries

```http
GET /api/deliveries
```

---

# Header Filtering

Before forwarding requests, the relay removes unsafe and hop-by-hop headers such as:

- Host
- Content-Length
- Transfer-Encoding
- Connection
- Keep-Alive

This prevents forwarding invalid transport-level metadata to downstream services.

---

# What I Learned

This project helped me better understand:

- HTTP request/response lifecycle
- Webhook architecture
- Request forwarding behavior
- Hop-by-hop headers
- JSON serialization
- JWT authentication
- PostgreSQL persistence
- Database migrations
- sqlc-generated queries
- UUID-based routing
- Delivery tracking patterns
- Frontend/backend integration
- Deployment workflows using Railway

---

# Post-MVP Roadmap

## Authentication and session management

- Token refresh and revocation flows
- Refresh token rotation
- Session expiration management
- Logout and token validation
- Persistent session tracking

## Delivery system

- Background worker queue
- Retry logic with backoff
- Delivery attempt limits
- Dead-letter handling for failed deliveries that exceed retry limits, allowing problematic events<br>to be isolated and inspected without blocking the rest of the system
- Async delivery processing to improve responsiveness<br>
  and ensure webhook deliveries do not block API requests

## Security

- Webhook signature verification to ensure incoming webhook requests are genuinely sent<br>by the platform and have not been tampered with during transit
- Endpoint secrets that allow each webhook consumer to authenticate requests<br>using a shared secret unique to their endpoint
- Request validation improvements to strengthen payload validation, enforce schema consistency,<br>and reduce the risk of malformed or malicious requests

## Scalability

- Concurrency control to manage how many delivery jobs run at the same time,<br>preventing system overload and ensuring stable performance under heavy traffic
- Rate limiting to restrict how frequently clients or endpoints can send requests,<br>protecting the service from abuse and maintaining reliability
- Delivery worker pools to process webhook deliveries asynchronously using multiple workers,<br>improving efficiency and reducing delivery delays
- Improved request throughput through performance optimizations that allow the system to handle<br>a larger number of requests and deliveries simultaneously

## Observability

- Structured logs to make debugging and tracing requests easier across services
- Metrics collection for tracking delivery success rates, latency, throughput, and system health
- Monitoring and alerting to detect failures, downtime, or performance degradation in real time
- Improved error tracing to better identify where and why delivery failures occur
- Delivery diagnostics tools for inspecting webhook attempts, responses, retries, and failure details

## API improvements

- Event filtering by type
- Pagination
- Search and filtering events, endpoints, or delivery attempts
- Better delivery inspection tools
- Add delivery retry endpoints that allow failed webhook deliveries to be manually retried
---

# Deployment

The service is deployed using:

- Railway
- PostgreSQL
- Goose migrations

The deployment process included:

- Environment variable configuration
- Production database migrations
- Public API routing
- Frontend/backend integration
- Git history recovery and safe force-pushing

---

# Contributing

Contributions, suggestions, and bug reports are welcome.

If you would like to contribute:

1. Fork the repository
2. Create a feature branch

```bash
git checkout -b feature/my-feature
```

3. Make your changes
4. Run tests and verify the application works correctly
5. Commit your changes

```bash
git commit -m "feat: add my feature"
```

6. Push your branch

```bash
git push origin feature/my-feature
```

7. Open a Pull Request

## Development Guidelines

- Follow standard Go formatting:

```bash
go fmt ./...
```
- Keep handlers focused on HTTP concerns
- Keep business logic inside service packages
- Add migrations for database schema changes
- Regenerate sqlc code after query updates:

```bash
sqlc generate
```

- Include tests when appropriate

## Reporting Issues

If you find a bug or have a feature request, please open an issue describing:

- Expected behavior
- Actual behavior
- Steps to reproduce
- Relevant logs or screenshots

All constructive feedback is appreciated.

---

# License

MIT