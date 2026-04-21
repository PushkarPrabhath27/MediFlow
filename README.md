# MediFlow – Hospital Equipment Utilisation & Sharing Optimiser

![MediFlow Architecture](https://img.shields.io/badge/Architecture-Go%20%7C%20Next.js%20%7C%20PostgreSQL%20%7C%20Redis-blue)
![License](https://img.shields.io/badge/License-MIT-green)

**MediFlow** is a real-time SaaS platform designed to solve a critical operational challenge in hospitals: the inefficient utilization and poor visibility of expensive medical equipment. By providing a centralized, real-time tracking system, MediFlow ensures that life-saving devices—such as ventilators, infusion pumps, and portable ultrasounds—are precisely where they need to be, exactly when they are needed.

## 🚀 The Problem We Solve
In a typical multi-department hospital, equipment often sits idle in one ward while another ward desperately searches for it. The lack of a centralized visibility system leads to:
- Time wasted manually searching for equipment via phone calls.
- Inefficient allocation of capital, as hospitals over-purchase equipment to compensate for poor visibility.
- Sub-optimal patient care due to delays in acquiring necessary devices.

## 💡 The MediFlow Solution
MediFlow acts as the nervous system for hospital asset management. It provides:
- **Live Availability Board:** Real-time visibility into the status and location of all shared equipment.
- **Structured Sharing Workflow:** A 5-step digital handoff process (Request, Match, Approve, Handoff, Return) replacing ad-hoc phone coordination.
- **Real-Time WebSockets & Redis Pub/Sub:** Instant updates across all connected clients without page refreshes.
- **Utilization Analytics:** Deep insights into equipment idle time, demand forecasting, and procurement justification.
- **Smart Alerts:** Automated warnings when minimum stock levels are breached.

## 🛠️ Technology Stack
MediFlow is built for high concurrency, low latency, and enterprise reliability:
- **Frontend:** Next.js 14 (App Router), React, Tailwind CSS, shadcn/ui
- **Backend:** Go (Golang) with Chi router
- **Real-Time Hub:** Gorilla WebSockets & Redis Pub/Sub
- **Database:** PostgreSQL 15
- **Caching & Queues:** Redis 7
- **Infrastructure:** Docker, Docker Compose, Kubernetes (AKS)

## 📂 Repository Structure
- `/backend`: The Go-based API gateway, equipment services, and WebSocket hub.
- `/frontend`: The Next.js client interface.
- `docker-compose.yml`: Local development environment setup containing Postgres, Redis, and application services.

## ⚙️ Quick Start

1. **Clone the repository**
   ```bash
   git clone <your-repository-url>
   cd mediflow
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Update the .env file with your local configurations
   ```

3. **Run the application via Docker Compose**
   ```bash
   docker-compose up --build
   ```

## 📊 Core Features At A Glance
- **Multi-tenant Architecture:** Isolated environments for different hospital tenants.
- **QR Code Integration:** Quick scan-to-update functionality for physical equipment tracking.
- **Demand Forecasting:** Algorithmic predictions for peak demand periods based on historical utilization logs.
