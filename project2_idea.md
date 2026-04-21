# MediFlow – Hospital Equipment Utilisation & Sharing Optimiser

## 1. Project Overview

**MediFlow** is a real-time SaaS platform that solves a critical and largely invisible problem in hospitals: expensive medical equipment sitting idle in one ward while another ward is desperately short of the same device. MediFlow gives hospitals complete visibility into equipment availability across all departments, enables inter-department sharing requests, tracks utilisation patterns, and produces actionable insights that help hospitals maximise the value of their equipment investment.

This is a problem that exists in every multi-department hospital in India, and no commercially available product addresses it well.

---

## 2. The Problem

### 2.1 The Core Scenario
Picture this:
- The ICU has two ventilators. One is in use. One has been idle for 3 days.
- The post-op ward urgently needs a ventilator tonight. Their only unit is down for maintenance.
- The post-op charge nurse calls around the hospital on the phone for 45 minutes trying to locate an available ventilator.
- Eventually they find one, push it across the building on a trolley.
- Total: 1 hour of nurse time wasted, 45 minutes of patient waiting.

This scenario plays out dozens of times daily across a large hospital — with ventilators, infusion pumps, ECG machines, portable ultrasounds, defibrillators, wheelchairs, and more.

### 2.2 Why It Happens
- **No centralised visibility:** Departments don't know what other departments have or whether it's in use.
- **Phone-based coordination:** Nurses and staff call each other to find available equipment — slow, unreliable, undocumented.
- **No utilisation data:** Hospital management has no data on which devices sit idle most of the time, making procurement decisions blind.
- **Equipment hoarding:** Departments over-request and stockpile equipment defensively because they don't trust availability.
- **No sharing workflow:** Even when equipment is shared, there's no formal handoff — no accountability for who has what, where.

### 2.3 Business Impact
- Hospitals over-purchase equipment to compensate for poor utilisation (wasted capital).
- Patient care delays caused by equipment location searches.
- Staff time wasted on manual coordination.
- Equipment goes missing or is not returned because there's no tracking.
- Finance has no data to justify equipment investments or identify redundant purchases.

### 2.4 Scale
- A 500-bed hospital might have 300–1,000 shared movable equipment items.
- Industry studies show average equipment utilisation in hospitals is 42–55% — meaning nearly half of it is idle at any given time.
- Improving utilisation by just 20% can defer ₹2–5 crore in unnecessary capital purchases annually.

---

## 3. The Solution

MediFlow provides:
1. **A live equipment availability board** — every shared device's current status: available, in use, in maintenance, in transit, reserved.
2. **A sharing request system** — a department can request equipment from another department through a structured digital workflow.
3. **Real-time status updates** via WebSockets — the board updates live without page refresh.
4. **Utilisation analytics** — which devices are used most, which sit idle, utilisation by department, time-of-day patterns.
5. **Equipment location tracking** — where is each device physically right now (room/ward level).
6. **Shortage alerts** — when a device category drops below minimum stock level in a department, alert goes out.
7. **Demand forecasting** — predict when peak demand periods are for specific equipment types.

---

## 4. Who Uses This (User Personas)

| Persona | Role | What They Get |
|---|---|---|
| Charge Nurse / Ward Sister | Manages ward equipment | Request equipment, check availability, see incoming/outgoing |
| Biomedical Engineer | Maintains equipment | Real-time location of all devices, maintenance scheduling |
| Department Head | Manages department resources | Department utilisation dashboard, sharing history |
| Hospital Administrator | Facility-wide oversight | Cross-department utilisation, procurement insights, cost savings report |
| Procurement Officer | Buys equipment | Utilisation data to justify or defer purchases |

---

## 5. Core Features

### 5.1 Equipment Catalogue & Status Registry
- Register all shared/movable equipment items: ventilators, ECG machines, infusion pumps, ultrasound machines, defibrillators, wheelchairs, etc.
- Each item has a status: Available, In Use, Reserved, In Maintenance, In Transit, Missing.
- Status updated manually via dashboard, or automatically via request workflow completions.
- Each item has a current location: Department + Room/Bay.
- QR code per device — staff scan to quickly update status.

### 5.2 Live Availability Board
- Real-time grid view: rows = departments, columns = equipment categories.
- Each cell shows count of available units in that department.
- Colour coding: green (available), amber (low — only 1 left), red (none available).
- Click any cell to drill down: see individual device names, status, current holder.
- Updates pushed via WebSocket — no page refresh needed.
- Filter by equipment category, building, floor.

### 5.3 Equipment Sharing Request Workflow
Structured 5-step workflow:

**Step 1 — Request:** Requesting department selects equipment category, quantity needed, urgency level, duration needed, reason.

**Step 2 — System matching:** MediFlow finds the best available equipment (closest location, lowest utilisation, not reserved) and suggests the source department.

**Step 3 — Approval:** Source department head or charge nurse approves or declines. If declined, system finds next best option.

**Step 4 — Handoff:** When physical handoff happens, both sides confirm via dashboard or QR scan. Equipment status changes to In Transit → In Use at requesting department.

**Step 5 — Return:** Requesting department marks as returned. Equipment returns to source department's Available pool.

Full audit trail for every sharing event.

### 5.4 Real-Time Notifications
- Alert when request is approved/declined.
- Alert when equipment that was requested (and unavailable) becomes available.
- Alert when a department's equipment stock in a category drops below their set minimum.
- Alert when equipment has been "In Transit" for more than a set time (potential loss).
- Delivered via Redis pub/sub → WebSocket to connected clients + in-app notification centre.

### 5.5 Utilisation Analytics
- **Utilisation Rate per device:** % of time each device is In Use vs Available over a time period.
- **Department utilisation heatmap:** Time-of-day × day-of-week grid showing when each department needs most equipment.
- **Equipment category demand:** Which categories are most requested across the hospital.
- **Sharing network graph:** Which departments share most with each other (visual).
- **Idle time tracker:** Devices with highest idle time — candidates for redeployment or disposal.
- **Procurement justification report:** Equipment categories that are consistently at 90%+ utilisation — clear case to buy more.
- **Cost savings report:** Estimated savings from sharing vs purchasing additional units.

### 5.6 Minimum Stock Level Alerts
- Each department can set a minimum stock level per equipment category.
- If current available count drops below minimum, alert fires.
- System proactively suggests borrowing from departments with excess.

### 5.7 Equipment Location History
- Full timeline of every location change for each device.
- "Where has this device been in the last 30 days?" — fully answerable.
- Useful for infection control (equipment that was in an isolation ward).

### 5.8 Multi-Tenancy
- Each hospital is a tenant.
- Fully isolated data.
- Subscription plans based on device count and department count.

---

## 6. Technical Architecture

### 6.1 High-Level Architecture

```
┌──────────────────────────────────────────────────────────┐
│                     CLIENT LAYER                         │
│           NextJS Frontend (SSR + CSR)                    │
│   Live Board | Requests | Analytics | Settings           │
│              WebSocket Connection (real-time)            │
└────────────────────────┬─────────────────────────────────┘
                         │ HTTPS / REST / WebSocket
┌────────────────────────▼─────────────────────────────────┐
│              API GATEWAY / ROUTER (GoLang)               │
│      JWT Auth | Rate Limiting | WS Upgrade Handler       │
└──────┬───────────┬────────────┬────────────┬─────────────┘
       │           │            │            │
┌──────▼──┐  ┌─────▼────┐  ┌───▼────┐  ┌────▼────────┐
│Equipment│  │ Request  │  │ Alert  │  │  Analytics  │
│ Service │  │ Service  │  │ Service│  │  Service    │
└─────────┘  └──────────┘  └────────┘  └─────────────┘
       │           │            │            │
┌──────▼───────────▼────────────▼────────────▼────────────┐
│               PostgreSQL (Primary Store)                 │
│  equipment | sharing_requests | utilisation_logs         │
│  equipment_locations | alerts | tenants | users          │
└──────────────────────────────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│                       Redis                              │
│  Availability Cache | WS Pub/Sub | Session Store         │
│  Real-time board state | Notification queues             │
└──────────────────────────────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│              WebSocket Hub (GoLang goroutines)           │
│  Broadcast availability updates to all connected clients │
└──────────────────────────────────────────────────────────┘
```

### 6.2 WebSocket Architecture (Key Innovation)
The live availability board is the heart of MediFlow. It uses:
- GoLang WebSocket hub (gorilla/websocket) with a broadcaster goroutine
- When any equipment status changes → save to DB → update Redis availability state → publish event to Redis channel → WebSocket hub receives event → broadcasts to all connected clients in that tenant
- Clients receive the delta update and update only the changed cell in the board UI

### 6.3 Database Schema (Key Tables)

```sql
equipment_items (id, tenant_id, name, category, serial_no, department_id, 
                 current_location_id, status, qr_code, purchase_date, created_at)

equipment_categories (id, tenant_id, name, icon, description)

locations (id, tenant_id, department_id, name, floor, building, type)

equipment_status_logs (id, equipment_id, old_status, new_status, location_id, 
                       changed_by_user_id, reason, changed_at)

sharing_requests (id, tenant_id, requesting_dept_id, source_dept_id, equipment_id, 
                  category_id, quantity_needed, urgency, reason, status, 
                  needed_by, expected_return_date, requested_at, approved_at, 
                  handed_off_at, returned_at, approved_by_user_id)

utilisation_logs (id, equipment_id, tenant_id, status, department_id, 
                  location_id, started_at, ended_at, duration_minutes)

department_min_stock (id, tenant_id, department_id, category_id, minimum_count)

demand_forecasts (id, tenant_id, category_id, department_id, forecast_date, 
                  predicted_demand, confidence, created_at)
```

### 6.4 Tech Stack

| Layer | Technology |
|---|---|
| Backend | GoLang (chi router + gorilla/websocket) |
| Frontend | NextJS 14 (App Router, SSR) |
| UI | React + TailwindCSS + shadcn/ui |
| Real-time | WebSockets (gorilla/websocket) + Redis Pub/Sub |
| Primary Database | PostgreSQL 15 |
| Caching | Redis 7 |
| Containerisation | Docker + Docker Compose |
| Orchestration | Kubernetes |
| CI/CD | GitHub Actions |
| Cloud | Azure AKS |
| API Docs | Swagger/OpenAPI 3.0 |
| Testing | Go testing + testify + k6 |

---

## 7. What Makes This Stand Out

1. **Real-time live board with WebSockets** — a technically impressive feature that feels genuinely useful.
2. **Solves an underserved niche** — no well-known commercial product targets this specific problem in India.
3. **Sharing workflow** — a proper state machine with approval, handoff, and return tracking.
4. **Utilisation analytics** — turns operational data into procurement decisions worth crores.
5. **Demand forecasting** — proactively predicts when equipment will be needed.
6. **QR code integration** — practical feature showing product thinking, not just backend skill.
7. **Redis pub/sub → WebSocket pipeline** — demonstrates understanding of real-time architecture.

---

## 8. Metrics to Demonstrate

- WebSocket connection handles 500 concurrent clients per pod.
- Availability board update latency: < 50ms from status change to all clients updated.
- Utilisation report for 1,000 devices over 6 months: generated in < 1 second (with proper indexing).
- Load test: 2,000 concurrent REST + WebSocket connections, p95 < 100ms.

---

## 9. Development Phases

| Phase | Deliverable | Duration |
|---|---|---|
| Phase 1 | Auth, tenant, equipment catalogue, locations | Week 1 |
| Phase 2 | Status management, live board, WebSocket hub | Week 2 |
| Phase 3 | Sharing request workflow, notifications | Week 3 |
| Phase 4 | Utilisation tracking, analytics, demand forecasting | Week 4 |
| Phase 5 | NextJS frontend (live board, requests, analytics) | Week 5 |
| Phase 6 | DevOps, testing, load testing, documentation | Week 6 |
