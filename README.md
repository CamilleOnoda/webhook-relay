# Webhook Relay Service

<p align="center">
  <img src="assets/system-overview.png" width="500">
</p>
An automated post office for HTTP messages. The relay receives and persists webhook events, delivers them asynchronously, and automatically retries transient failures using exponential backoff with jitter. Delivery state is persisted across restarts, exhausted deliveries are moved to a dead-letter queue, and delivery activity can be inspected and replayed through user and admin dashboards.

Built in Go with PostgreSQL, with a focus on reliable event delivery, failure handling, observability, and operating and troubleshooting the service in a Linux environment.

---

## Repository highlights

- **Asynchronous webhook delivery** — incoming webhooks are accepted quickly, while delivery processing happens in the background.
- **Automatic retries with exponential backoff** — failed deliveries are rescheduled instead of being lost immediately.
- **JWT + refresh token authentication** — users stay authenticated with short-lived access tokens and secure refresh cookies.
- **User & admin monitoring dashboards** — delivery status, retries, dead letters, endpoints, and users can be inspected from the UI.

---

# Live demo
Try the deployed application:

👉 https://webhook-relay-production-5e97.up.railway.app/

### Demo credentials

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

# Table of contents

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
- [Next engineering focus](#next-engineering-focus)
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

### Delivery lifecycle

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

## Webhook management
- Create and manage webhook endpoints
- Generate unique webhook URLs
- Store endpoint metadata in PostgreSQL
- Delete endpoints and associated data

## Event @rocessing
- Receive incoming webhook events
- Store payloads and metadata
- Filter unsafe transport headers
- Forward requests to destination services
- Track delivery status and response codes

## Reliable delivery
- Automatic retry scheduling
- Exponential backoff
- Retry jitter to prevent retry storms
- Delivery attempt tracking
- Dead-letter queue support
- Dead-letter replay
- Worker recovery after restarts

## Authentication & authorization
- User registration and login
- Argon2id password hashing
- JWT access tokens
- Refresh token cookies
- Automatic token refresh
- Session revocation
- Route protection middleware
- Admin-only endpoints

## Dashboard

### User dashboard
- Endpoint management
- Event history
- Delivery history
- Successful delivery metrics
- Retry scheduled metrics
- Dead-letter metrics
- Delivery status tracking

### Admin dashboard
- System-wide statistics
- User management
- Endpoint overview
- Recent activity feed
- Dead-letter inspection
- Delivery monitoring
- Delivery details modal
- Dead-letter replay
- Attempt count monitoring

## Implementation details
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

## Delivery lifecycle
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

## Session flow

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
      <img src="assets/session-flow.png" width="700">
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

# Project structure

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

# Quick start
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
## Create an account
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

## Create a webhook endpoint
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

## Send a test webhook
There are two ways to generate webhook events:

- Click **Send Test** from the User Dashboard.
- Send a POST request directly to your webhook URL:

```bash
curl -X POST http://localhost:8080/webhooks/{endpoint_id} \
-H "Content-Type: application/json" \
-d '{"type":"payment.success"}'
```

## Inspect results
```
GET /api/events
GET /api/deliveries
```

## View:

- Stored events
- Delivery attempts
- Response status codes
- Retry status

## Admin dashboard
Administrators can access:

- System statistics
- User management
- Endpoint monitoring
- Recent activity
- Dead-letter queue inspection

### Admin API routes
```
GET /admin/stats
GET /admin/users
GET /admin/endpoints
GET /admin/recent-activity
```

---

# What I learned

The biggest lesson from this project was that reliability changes the architecture of an application.

The first version of the relay could have simply received a webhook and forwarded it immediately. Once I wanted deliveries to survive temporary failures, however, the problem became much more interesting.

Events needed to be persisted before delivery. Delivery attempts needed their own state. Failed requests had to be scheduled for later rather than retried immediately, and retry information had to survive application restarts.

This led me to work through questions such as:

- What happens to a delivery if the application stops after receiving the event?
- Which state needs to be persisted rather than kept in memory?
- How does a worker know which deliveries are ready to be attempted?
- How should temporary failures differ from permanently failed deliveries?
- How can retries avoid repeatedly hammering an unhealthy destination?
- How can failed deliveries be replayed without creating a separate delivery path?

Working through those problems led to the current design: persisted events and delivery state, background workers, exponential backoff with jitter, dead-letter handling, and a replay mechanism that sends failed deliveries back through the normal delivery lifecycle.

The project also gave me practical experience with authentication and session management using short-lived JWT access tokens and refresh tokens, PostgreSQL-backed application state, database migrations, generated queries with sqlc, and deploying a Go application with a production database.

More importantly, it changed how I think about backend systems. A successful HTTP request is only one part of the problem; I also need to think about what happens when dependencies are slow, unavailable, or when the application itself restarts.

---

# Next engineering focus

The core delivery workflow is implemented, so the next phase of the project is less about adding application features and more about understanding how the service behaves while it is running.

**Observability**

- Introduce structured application logging
- Expose application and delivery metrics
- Collect metrics with Prometheus
- Visualize system behaviour with Grafana
- Track delivery success/failure rates, retries, dead letters, and worker activity

**Linux deployment and operations**

Run the service in an Ubuntu environment and use it as a small production-like system for practising:
- Process and service management
- Networking and port troubleshooting
- PostgreSQL connectivity
- Application logs
- Resource monitoring
- Service startup and recovery
- Failure diagnosis

**Failure testing**

Deliberately introduce failures and investigate them from the system rather than only from the application code.

Examples include:
- Destination endpoint unavailable
- Destination endpoint returning 5xx responses
- PostgreSQL unavailable
- Network or DNS failure
- Delivery worker stopped unexpectedly
- Application restart with pending deliveries
- Resource or disk pressure

The goal is to observe how failures appear across logs, metrics, processes, networking, the database, and the delivery lifecycle, then diagnose the problem from those signals.

**Security**

- Webhook signatures
- Endpoint secrets
- Stronger request validation

**Possible future architecture experiments**

As the project grows, I may also use it to experiment with alternative approaches such as:

- Queue-based delivery processing
- Worker pools
- Concurrency controls
- Rate limiting

These are potential architectural experiments rather than requirements for the current implementation.

---

# Contributing

Contributions are very welcome!

This project started as a way for me to explore reliable webhook delivery, but I'd love for it to become something other people can experiment with, improve, and learn from too.

There are many ways to contribute, including:

- Fixing bugs
- Improving tests
- Improving documentation
- Improving logging and observability
- Suggesting or investigating reliability edge cases
- Improving development or deployment tooling
- Opening an issue with an idea or unexpected behaviour you found

You don't need to take on a large feature. Small fixes, tests, documentation improvements, and discussions are just as welcome.

If you're interested in contributing but aren't sure where to start, feel free to open an issue or look through the existing issues.

**Making a Contribution**

1. Fork the repository and create a branch:

`git switch -c fix/my-change`

2. Make your changes.
3. Run the tests and verify the application still works as expected.
4. Format Go code:

`go fmt ./...`

5. Commit and push your changes, then open a Pull Request describing what you changed and why.

**Development Guidelines**

A few guidelines to keep the codebase consistent:

- Keep handlers focused on HTTP concerns.
- Keep business logic inside service packages.
- Add migrations for database schema changes.
- Regenerate sqlc code after modifying queries:

`sqlc generate`

- Add or update tests when the change affects existing behaviour.
- Keep changes focused. A small PR that solves one problem is easier to review than several unrelated changes bundled together.

If you're considering a larger architectural change, opening an issue first is encouraged so we can discuss the approach before you spend time implementing it.

**Reporting bugs and suggesting improvements**

Issues are welcome too.

For bugs, include whatever information you have that could help reproduce or understand the problem:

- What you expected to happen
- What actually happened
- Steps to reproduce it
- Relevant logs or screenshots
- Your environment, when relevant

Ideas, questions, reliability edge cases, and suggestions for improving the project are also welcome. You don't need to already have a solution before opening an issue.

---

# License

MIT