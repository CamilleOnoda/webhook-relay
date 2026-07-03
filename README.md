# Webhook Relay Service

A webhook relay service written in Go that receives webhook events, stores them in PostgreSQL, delivers them to downstream services, and automatically retries failed deliveries using exponential backoff.

Built to explore reliable event delivery, authentication,
background processing, and backend system design in Go.

---

# Live Demo
Try the deployed application:

👉 https://webhook-relay-production-5e97.up.railway.app/

### Demo Credentials

If you'd like to explore the application without creating an account, use one of the demo users:

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@test.com | password1234 |
| User | user@test.com | password1234 |

You can also register your own account if you'd like to test the registration flow.

---

# Motivation

Many modern applications rely on webhooks to communicate between services, yet the mechanics behind webhook delivery are often hidden behind third-party platforms.

I built this project to better understand what happens between the moment a webhook is received and the moment it reaches its destination.

The goal was not only to learn Go, but to gain experience with:

- Receiving incoming requests
- Validating and storing data
- Forwarding events to destination services
- Handling delivery failures
- Retrying failed deliveries
- Managing user authentication and sessions
- Monitoring system activity

Rather than consuming webhooks through existing tools, I wanted to build the infrastructure myself and explore the challenges involved in receiving, storing, processing, and forwarding webhook events reliably.

This project became an opportunity to combine several backend concepts into a single application while simulating the kind of event-driven systems commonly used in production environments.

---

# Table of Contents

- [Motivation](#motivation)
- [Screenshots](#screenshots)
- [Features](#features)
- [Architecture](#architecture)
- [Session Flow](#session-flow)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [What I Learned](#what-i-learned)
- [Future Improvements](#future-improvements)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Screenshots

### Dashboards

<table>
  <tr>
    <td align="center" width="33%">
      <strong>User Dashboard</strong><br>
      Manage webhook endpoints, monitor delivery metrics, and inspect recent activity.
      <br><br>
      <img src="assets/user-dashboard.png" width="400">
    </td>
    <td align="center" width="33%">
      <strong>Admin Dashboard</strong><br>
      Monitor users, endpoints, delivery outcomes, retry queues, and failed deliveries.
      <br><br>
      <img src="assets/admin-dashboard.png" width="400">
    </td>
    <td align="center" width="33%">
      <strong>Management Tables</strong><br>
      View registered users and webhook endpoints from the admin dashboard.
      <br><br>
      <img src="assets/endpoints-users.png" width="400">
    </td>
  </tr>
</table>

### Delivery Lifecycle

Successful deliveries, scheduled retries, and dead-lettered events can be inspected individually.

<table>
  <tr>
    <td align="center">
      <strong>Success</strong><br>
      <img src="assets/delivery-success.png" width="300">
    </td>
    <td align="center">
      <strong>Retry Scheduled</strong><br>
      <img src="assets/delivery-retry.png" width="300">
    </td>
    <td align="center">
      <strong>Dead Letter</strong><br>
      <img src="assets/delivery-deadletter.png" width="300">
    </td>
  </tr>
</table>

---

# Features

## Webhook Management
- Create and manage webhook endpoints
- Generate unique webhook URLs
- Store endpoint metadata in PostgreSQL
- Delete endpoints and associated data

## Event Processing
- Receive incoming webhook events
- Store payloads and metadata
- Filter unsafe transport headers
- Forward requests to destination services
- Track delivery status and response codes

## Reliable Delivery
- Automatic retry scheduling
- Exponential backoff
- Retry jitter to prevent retry storms
- Delivery attempt tracking
- Dead-letter queue support
- Dead-letter replay
- Worker recovery after restarts

## Authentication & Authorization
- User registration and login
- Argon2id password hashing
- JWT access tokens
- Refresh token cookies
- Automatic token refresh
- Session revocation
- Route protection middleware
- Admin-only endpoints

## Dashboard

### User Dashboard
- Endpoint management
- Event history
- Delivery history
- Successful delivery metrics
- Retry scheduled metrics
- Dead-letter metrics
- Delivery status tracking

### Admin Dashboard
- System-wide statistics
- User management
- Endpoint overview
- Recent activity feed
- Dead-letter inspection
- Delivery monitoring
- Delivery details modal
- Dead-letter replay
- Attempt count monitoring

## Implementation Details
### Header filtering
Before forwarding requests, the relay removes hop-by-hop headers such as:

- Host
- Content-Length
- Connection
- Transfer-Encoding
- Keep-Alive

This prevents transport-level metadata from being incorrectly forwarded to downstream services.

---

## Architecture

Webhook delivery is intentionally separated from request reception.
When a webhook is received:

1. Endpoint is validated
2. Event is persisted
3. Delivery record is created
4. Delivery worker processes the event
5. Success or failure is recorded
6. Failed deliveries are rescheduled
7. Exhausted deliveries move to the dead-letter queue
8. Administrators can replay dead-letter deliveries
9. Replayed deliveries re-enter the delivery workflow

This design allows delivery processing to continue independently of incoming requests and provides a foundation for future queue-based processing.

<p align="center">
  <img src="assets/webhook-architecture.png" width="500">
</p>

## Delivery Lifecycle
```text
pending
    ↓
success
```

```text
pending
    ↓
retry_scheduled
    ↓
retry_scheduled
    ↓
dead_letter
    ↓
admin replay
    ↓
pending
```

## Session Flow

Login creates two credentials:

- Short-lived JWT access token
- Long-lived refresh token stored as an HttpOnly cookie

Flow:

Login
→ Access Token + Refresh Cookie
→ Protected Request
→ Access Token Expires
→ Refresh Endpoint
→ New Access Token
→ Retry Original Request

This allows sessions to remain active without repeatedly prompting users to log in.

<table>
  <tr>
    <td align="center">
      <img src="assets/session-flow.png" width="500">
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
├── assets/
├── internal/
│   ├── auth/
│   ├── database/
│   ├── service/
│   │   ├── delivery.go
│   │   ├── delivery_integration_test.go
│   │   ├── helper_unsafeHeaders.go
│   │   └── is_retryable.go
│   └── static/
│       ├── admin.html
│       ├── app.js
│       ├── dashboard.html
│       ├── index.html
│       └── styles.css
├── sql/
│   ├── queries/
│   │   ├── admin.sql
│   │   ├── auth.sql
│   │   ├── deliveries.sql
│   │   ├── endpoints.sql
│   │   ├── events.sql
│   │   ├── refresh_tokens.sql
│   │   └── users.sql
│   └── schema/
├── admin_middleware.go
├── auth_middleware.go
├── create_endpoint_test.go
├── handler_admin_dead_letter_get.go
├── handler_admin_dead_letter_replay.go
├── handler_admin_deliveries_get.go
├── handler_admin_endpoints_get.go
├── handler_admin_events_get.go
├── handler_admin_recent_activity.go
├── handler_admin_stats_get.go
├── handler_admin_user_delete.go
├── handler_admin_users_get.go
├── handler_deliveries_get.go
├── handler_endpoint_create.go
├── handler_endpointByID_delete.go
├── handler_endpointByID_get.go
├── handler_endpoints_delete.go
├── handler_endpoints_get.go
├── handler_events_get.go
├── handler_login.go
├── handler_token_refresh.go
├── handler_token_revoke.go
├── handler_users_create.go
├── handler_webhook_receive.go
├── json.go
├── main.go
├── readiness.go
├── go.mod
├── go.sum
├── sqlc.yaml
├── test_REST_client.http
└── README.md
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
## Create an Account
```
POST /api/users
```

## Login
```
POST /api/login
```

Returns:
- JWT access token
- Refresh token cookie

## Create a Webhook Endpoint
```
POST /api/endpoints
```

Response:
```
{
  "id": "...",
  "generated_url": "/webhooks/{id}"
}
```

## Send a Test Webhook
There are two ways to generate webhook events:

- Click **Send Test** from the User Dashboard.
- Send a POST request directly to your webhook URL:

```bash
curl -X POST http://localhost:8080/webhooks/{endpoint_id} \
-H "Content-Type: application/json" \
-d '{"type":"payment.success"}'
```

## Inspect Results
```
GET /api/events
GET /api/deliveries
```

## View:

- Stored events
- Delivery attempts
- Response status codes
- Retry status

## Admin Dashboard
Administrators can access:

- System statistics
- User management
- Endpoint monitoring
- Recent activity
- Dead-letter queue inspection

### Admin API Routes
```
GET /admin/stats
GET /admin/users
GET /admin/endpoints
GET /admin/recent-activity
```

---

# What I Learned

Building this project helped me gain hands-on experience with:

- REST API design
- Authentication and authorization
- Session management
- JWTs and refresh tokens
- PostgreSQL schema design
- Database migrations
- sqlc
- Webhook delivery systems
- Retry and backoff strategies
- Dead-letter queues
- Background workers
- Frontend/backend integration
- Railway deployments
- Production debugging
- Background job processing
- Failure recovery strategies
- Dead-letter queue design
- Operational monitoring dashboards

---

# Future Improvements

## Authentication and session management

- Refresh token rotation
- Multiple concurrent sessions
- Logout from all devices
- Session management dashboard

## Security

- Webhook signatures
- Endpoint secrets
- Enhanced request validation

## Scalability

- Queue-based delivery processing
- Worker pools
- Rate limiting
- Concurrency controls

## Observability

- Structured logging
- Metrics
- Alerting
- Delivery analytics

## API improvements

- Event filtering by type
- Pagination
- Search and filtering events, endpoints, or delivery attempts

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