# MediFlow

MediFlow is a full-stack hospital equipment utilization and sharing platform built to improve visibility, reduce idle inventory, and coordinate inter-department equipment movement with a structured digital workflow.

## Overview

This repository contains:

- A `Next.js` frontend for availability, request workflow, and dashboard views
- A `Go` backend exposing REST APIs, JWT-based auth, analytics, equipment management, and request orchestration
- `PostgreSQL` for persistent multi-tenant operational data
- `Redis` for real-time coordination and WebSocket event propagation
- `Docker Compose` for local development orchestration

## Key Capabilities

- Real-time equipment visibility across departments
- Structured equipment sharing request lifecycle
- QR-based equipment status updates
- Availability summaries and status history
- Dashboard analytics for operational monitoring
- Multi-tenant data separation
- WebSocket updates backed by Redis Pub/Sub

## Tech Stack

- Frontend: `Next.js`, `React`, `TypeScript`, `Tailwind CSS`
- Backend: `Go`, `chi`, `sqlx`, `zerolog`
- Data: `PostgreSQL`, `Redis`
- Infrastructure: `Docker`, `Docker Compose`

## Repository Structure

```text
.
|-- backend/        Go API, domain services, migrations, and WebSocket hub
|-- frontend/       Next.js application
|-- docker-compose.yml
|-- .env.example
|-- project2_*.md   product and planning artifacts
```

## Getting Started

### Prerequisites

- `Docker` and `Docker Compose`

### Local Setup

1. Clone the repository.
2. Copy the example environment file:

```bash
cp .env.example .env
```

3. Start the full stack:

```bash
docker compose up --build
```

### Default Local Endpoints

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- Health check: `http://localhost:8080/health`

## Environment Configuration

Core variables are defined in [.env.example](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/.env.example>).

Important values include:

- `PORT`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `REDIS_HOST`
- `REDIS_PORT`
- `JWT_SECRET`
- `FRONTEND_URL`

Never commit real secrets or production credentials.

## Backend Notes

The backend entry point is [backend/cmd/server/main.go](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/backend/cmd/server/main.go>). It initializes:

- database and Redis connections
- JWT authentication
- equipment, request, alert, and analytics services
- Redis-backed WebSocket broadcasting
- background job workers

Primary API groups currently include:

- `/api/v1/auth`
- `/api/v1/equipment`
- `/api/v1/requests`
- `/api/v1/analytics`
- `/api/v1/ws`

## Frontend Notes

The frontend lives in [frontend](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/frontend>) and includes:

- dashboard overview page
- live availability board
- sharing workflow view

## Quality and Maintenance

- Use `.env.example` as the source of truth for local configuration
- Keep documentation aligned with implemented behavior
- Prefer small, reviewable pull requests
- Add or update tests when behavior changes

## Documentation

- [CONTRIBUTING.md](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/CONTRIBUTING.md>)
- [SECURITY.md](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/SECURITY.md>)
- [CODE_OF_CONDUCT.md](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/CODE_OF_CONDUCT.md>)
- [frontend/README.md](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/frontend/README.md>)

## License

This project is released under the MIT License. See [LICENSE](</C:/Users/pushk/OneDrive/Documents/MediFlow 3/LICENSE>).
