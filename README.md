# Webhook Relay Service

<p align="center">
  <img src="assets/system-overview.png" width="600">
</p>
An automated post office for HTTP messages. The relay receives webhook events, stores them, delivers them asynchronously, automatically retries transient failures using exponential backoff, and exposes delivery status through user and admin dashboards. Built in Go to explore reliable event delivery, authentication, background workers, and backend system design.

---

## Repository Highlights

- 🚀 **Asynchronous webhook delivery** — incoming webhooks are accepted quickly, while delivery processing happens in the background.
- 🔄 **Automatic retries with exponential backoff** — failed deliveries are rescheduled instead of being lost immediately.
- 🔐 **JWT + refresh token authentication** — users stay authenticated with short-lived access tokens and secure refresh cookies.
- 📊 **User & admin monitoring dashboards** — delivery status, retries, dead letters, endpoints, and users can be inspected from the UI.

---

# Live Demo
Try the deployed application:

👉 https://webhook-relay-production-5e97.up.railway.app/

### Demo Credentials

Use the demo credentials below to explore the application safely.
The admin dashboard displays registered user emails, so please avoid using a personal email address when testing.

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@demo.com | password1234 |
| User | user@demo.com | password1234 |

You can also register your own account if you'd like to test the registration flow.

---

# Motivation

Many modern applications rely on webhooks, but I realized I had never really thought about what happened after an external service sent an event. I wanted to build that missing piece myself.

I started thinking of a webhook relay as an automated post office for HTTP messages—receiving events, keeping track of deliveries, retrying failed ones, and making sure messages eventually reach their destination.
That simple idea became an opportunity to explore authentication, background workers, retry scheduling, dead-letter queues, and reliable event delivery in Go.

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

```mermaid
flowchart TD
    Client[Client Application]
    API[Webhook Relay API]
    DB[(PostgreSQL)]
    Worker[Background Delivery Worker]
    Target[Destination Endpoint]

    Client -->|POST /webhooks/:id| API

    API -->|Persist Event| DB

    Worker -->|Fetch Pending Deliveries| DB
    Worker -->|Deliver Webhook| Target
    Worker -->|Record Success / Failure| DB
```

## Delivery Lifecycle
```mermaid
flowchart TD
    A[Webhook Received] --> B[Pending Delivery]

    B -->|Success| C[Delivered]

    B -->|Failure| D[Retry Scheduled]
    D --> E[Exponential Backoff]
    E --> F{Retries Left?}

    F -->|Yes| B
    F -->|No| G[Dead-letter Queue]

    G -->|Admin Replay| B
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
internal/
    auth/
    database/
    service/
    static/

sql/
    schema/
    queries/

assets/

main.go
README.md
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

**Backend**
- REST API design
- Authentication
- Session management

**Reliability**
- Retry strategies
- Dead-letter queues
- Background workers

**Infrastructure**
- Railway
- PostgreSQL
- Goose
- sqlc

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