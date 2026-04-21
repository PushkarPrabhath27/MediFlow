# MediFlow – AI Implementation Prompt (Phase-wise)

> **Instructions for use:** Feed each phase to your AI coding assistant (Cursor, Claude, Copilot) separately. Complete and verify each phase before moving to the next.

---

## GLOBAL CONTEXT (Include at the start of every phase prompt)

```
You are building MediFlow — a real-time multi-tenant SaaS platform for hospital equipment utilisation and sharing optimisation.

Tech stack:
- Backend: GoLang (net/http + chi router + sqlx + gorilla/websocket)
- Frontend: NextJS 14 (App Router)
- Database: PostgreSQL 15
- Cache/PubSub: Redis 7
- Real-time: WebSockets (gorilla/websocket) + Redis Pub/Sub
- Containerisation: Docker + Docker Compose
- Orchestration: Kubernetes
- CI/CD: GitHub Actions
- API Docs: Swagger/OpenAPI 3.0

Project structure:
mediflow/
├── backend/
│   ├── cmd/
│   │   └── server/main.go
│   ├── internal/
│   │   ├── auth/
│   │   ├── equipment/
│   │   ├── location/
│   │   ├── request/
│   │   ├── utilisation/
│   │   ├── analytics/
│   │   ├── alert/
│   │   ├── tenant/
│   │   └── shared/
│   │       ├── db/
│   │       ├── redis/
│   │       ├── websocket/
│   │       ├── middleware/
│   │       └── models/
│   ├── migrations/
│   ├── docs/
│   └── Dockerfile
├── frontend/
│   ├── app/
│   ├── components/
│   └── lib/
├── k8s/
├── .github/workflows/
└── docker-compose.yml

Always:
- Write idiomatic GoLang with proper error handling
- Use repository pattern for all database access
- Write table-driven unit tests for all business logic
- Add Swagger annotations to all handlers
- Use structured logging with zerolog
- All DB queries use parameterised statements
- Every response includes request_id
- WebSocket messages use typed JSON envelopes
```

---

## PHASE 1 — Foundation, Auth & Equipment Catalogue

### Prompt:

```
Using the global context above, implement Phase 1 of MediFlow.

TASK 1 — Project Scaffolding:
Create complete directory structure. Initialise:
- Go module: github.com/yourusername/mediflow
- Dependencies: chi, sqlx, lib/pq, go-redis, zerolog, golang-jwt/jwt/v5, 
  google/uuid, gorilla/websocket, swaggo/swag, testify
- NextJS with tailwindcss, shadcn/ui, axios, react-query, zustand

TASK 2 — Docker Compose:
docker-compose.yml with:
- PostgreSQL 15 (health check, persistent volume)
- Redis 7 (persistent volume, enable keyspace notifications: "KEA" — required for pub/sub)
- Backend with hot reload (air)
- Frontend with hot reload
- Shared network
- .env.example with all variables documented

TASK 3 — Database Migrations:

Migration 001 — tenants:
CREATE TABLE tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(100) UNIQUE NOT NULL,
  plan VARCHAR(50) NOT NULL DEFAULT 'basic',
  max_devices INTEGER NOT NULL DEFAULT 100,
  max_departments INTEGER NOT NULL DEFAULT 10,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 002 — users:
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  email VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  role VARCHAR(50) NOT NULL DEFAULT 'staff',
  department_id UUID,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, email)
);

Migration 003 — departments:
CREATE TABLE departments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  name VARCHAR(255) NOT NULL,
  code VARCHAR(20),
  floor VARCHAR(50),
  building VARCHAR(100),
  head_user_id UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 004 — locations:
CREATE TABLE locations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  name VARCHAR(255) NOT NULL,
  type VARCHAR(50), -- room, bay, corridor, storage
  floor VARCHAR(50),
  building VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 005 — equipment_categories:
CREATE TABLE equipment_categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  icon VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 006 — equipment_items:
CREATE TABLE equipment_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  department_id UUID REFERENCES departments(id),
  current_location_id UUID REFERENCES locations(id),
  name VARCHAR(255) NOT NULL,
  model VARCHAR(255),
  manufacturer VARCHAR(255),
  serial_number VARCHAR(255),
  asset_tag VARCHAR(100),
  qr_code VARCHAR(255) UNIQUE,
  status VARCHAR(50) NOT NULL DEFAULT 'available',
  -- status: available, in_use, reserved, in_maintenance, in_transit, missing, decommissioned
  purchase_date DATE,
  purchase_cost DECIMAL(12,2),
  notes TEXT,
  is_shared BOOLEAN NOT NULL DEFAULT true, -- whether this item participates in sharing
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_equipment_tenant_status ON equipment_items(tenant_id, status);
CREATE INDEX idx_equipment_department ON equipment_items(department_id);
CREATE INDEX idx_equipment_category ON equipment_items(category_id);

Migration 007 — department_min_stock:
CREATE TABLE department_min_stock (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  minimum_count INTEGER NOT NULL DEFAULT 1,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(department_id, category_id)
);

TASK 4 — DB & Redis Setup:
Identical pattern to Phase 1 of MediGuard — connection pool, health checks, retry logic, Redis helpers.
Additionally for Redis: implement a GetAvailabilityKey(tenantID, departmentID, categoryID) helper that returns a consistent Redis key pattern: availability:{tenantID}:{departmentID}:{categoryID}

TASK 5 — Auth Service:
Same pattern as MediGuard Phase 1 auth. Roles for MediFlow:
- super_admin: platform admin
- hospital_admin: full tenant access
- department_head: manage own department, approve/decline sharing requests
- charge_nurse: create requests, update equipment status
- staff: view only, update status for equipment they check out
- engineer: update maintenance status

TASK 6 — Equipment Catalogue Service:
Implement backend/internal/equipment/ with full repository/service/handler:

Repository:
- Create(item)
- FindByID(tenantID, itemID)
- FindAll(tenantID, filters, pagination) — filters: category_id, department_id, status, is_shared, search
- Update(item)
- UpdateStatus(itemID, newStatus, locationID, changedByUserID, reason)
- FindByQRCode(qrCode)
- CountAvailableByCategory(tenantID, departmentID, categoryID)
- FindAvailableByCategory(tenantID, categoryID) — for sharing suggestions
- GetAvailabilitySummary(tenantID) — aggregate counts per dept per category for board

Service:
- CreateItem: validate, generate QR code, check device limit
- UpdateStatus: validate transition, save status log, invalidate Redis availability cache, publish change event
- GetAvailabilitySummary: check Redis cache first (TTL 30s), if miss query DB and cache

Valid status transitions:
- available → in_use, reserved, in_maintenance, in_transit, missing
- in_use → available, in_maintenance, missing
- reserved → in_use, available (reservation cancelled)
- in_maintenance → available, decommissioned
- in_transit → available, in_use
- missing → available (found), decommissioned
Any other transition returns 400.

API Endpoints:
- POST /api/v1/equipment
- GET /api/v1/equipment (paginated, filterable)
- GET /api/v1/equipment/:id
- PUT /api/v1/equipment/:id
- DELETE /api/v1/equipment/:id
- PUT /api/v1/equipment/:id/status
- GET /api/v1/equipment/qr/:qrCode
- GET /api/v1/equipment/availability-summary (used by live board)
- POST /api/v1/equipment/bulk-import (CSV)

Also implement full CRUD for:
- Equipment categories: POST/GET/PUT/DELETE /api/v1/categories
- Locations: POST/GET/PUT/DELETE /api/v1/locations
- Departments: POST/GET/PUT/DELETE /api/v1/departments
- Min stock levels: POST/PUT /api/v1/departments/:id/min-stock

TASK 7 — Equipment Status Log:
Migration 008 — equipment_status_logs:
CREATE TABLE equipment_status_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  equipment_id UUID NOT NULL REFERENCES equipment_items(id),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50) NOT NULL,
  old_department_id UUID,
  new_department_id UUID,
  old_location_id UUID,
  new_location_id UUID,
  changed_by_user_id UUID REFERENCES users(id),
  reason TEXT,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_status_logs_equipment ON equipment_status_logs(equipment_id);
CREATE INDEX idx_status_logs_changed_at ON equipment_status_logs(changed_at);

Every call to UpdateStatus must insert a status log record. Expose:
- GET /api/v1/equipment/:id/history (full status history, paginated)
- GET /api/v1/equipment/:id/location-history (only location changes)

TASK 8 — Tests:
- Status transition state machine: all valid, all invalid
- QR code uniqueness enforcement
- Device limit enforcement per tenant plan
- Availability cache invalidation on status change
```

---

## PHASE 2 — WebSocket Hub & Live Availability Board

### Prompt:

```
Using the global context above, implement Phase 2 of MediFlow.
Phase 1 complete. Auth, equipment catalogue, status management work.

TASK 1 — WebSocket Hub:
Implement backend/internal/shared/websocket/hub.go:

Hub struct:
- clients: map[string]map[*Client]bool (key is tenantID — tenant isolation at hub level)
- register chan *Client
- unregister chan *Client  
- broadcast chan BroadcastMessage
- mu sync.RWMutex

Client struct:
- conn *websocket.Conn
- tenantID string
- userID string
- send chan []byte

BroadcastMessage struct:
- TenantID string
- Type string (availability_update, alert, request_update, notification)
- Payload interface{}

Hub.Run() goroutine:
- Loop on select:
  - register: add client to tenant bucket
  - unregister: remove client, close send channel, close connection
  - broadcast: for each message, only broadcast to clients in the matching tenantID bucket

Client.ReadPump() goroutine:
- Read messages from WebSocket connection
- Handle ping/pong for connection keepalive
- On disconnect: unregister from hub

Client.WritePump() goroutine:
- Write messages from send channel to WebSocket connection
- Handle write deadlines (5 second timeout per write)
- Ticker for ping messages every 30 seconds

WebSocket endpoint handler:
- GET /api/v1/ws — upgrades connection
- Requires valid JWT (passed as query param: ?token=... since WebSocket headers are limited)
- Extract tenantID and userID from token
- Create Client, register with hub, start ReadPump and WritePump goroutines

TASK 2 — Redis to WebSocket Bridge:
In backend/internal/shared/websocket/bridge.go:

Implement a StartBridge(hub *Hub, redisClient *redis.Client) function that:
1. Subscribes to Redis channel pattern: ws-events:* using PSubscribe
2. Runs as goroutine
3. On each message: deserialise JSON, extract tenantID, create BroadcastMessage, send to hub.broadcast channel
4. Log and continue on deserialisation errors
5. Handle subscription errors with reconnection (retry every 5s)

This bridge is the central mechanism: anything that publishes to ws-events:{tenantID} in Redis will be pushed to all WebSocket clients for that tenant.

TASK 3 — Availability State in Redis:
Create a service in backend/internal/equipment/ called AvailabilityStateManager:

On startup (and via a scheduled refresh every 60s):
- For all tenants, load the current availability summary from DB
- Store in Redis as a hash: key = availability_board:{tenantID}, field = {deptID}:{categoryID}, value = {available_count, in_use_count, total_count}
- TTL = 90 seconds (refreshed every 60s so it never expires during normal operation)

When any equipment status changes (called from UpdateStatus service method):
1. Re-query the specific dept+category combination from DB
2. Update only that field in the Redis hash
3. Publish to Redis channel: ws-events:{tenantID} with payload:
   {
     "type": "availability_update",
     "data": {
       "department_id": "...",
       "category_id": "...",
       "available_count": 3,
       "in_use_count": 2,
       "reserved_count": 0,
       "total_count": 5
     }
   }
4. The bridge picks this up and broadcasts to all connected clients

TASK 4 — Board State API:
GET /api/v1/availability-board
- Returns the full availability summary for the tenant
- Reads from Redis hash first (if available and fresh)
- Falls back to DB query
- Response format:
  {
    "board": [
      {
        "department": {id, name, floor},
        "categories": [
          {
            "category": {id, name, icon},
            "available": 3,
            "in_use": 2,
            "reserved": 0,
            "in_maintenance": 1,
            "total": 6,
            "below_minimum": false
          }
        ]
      }
    ],
    "last_updated": "ISO8601"
  }

GET /api/v1/availability-board/category/:categoryID
- Returns availability for a specific category across all departments
- Useful for "find me a device of this type" use case

GET /api/v1/departments/:id/equipment
- Returns all individual equipment items in a department
- Grouped by status
- Includes current location for each item

TASK 5 — QR Code Status Update Endpoint:
POST /api/v1/equipment/qr-update
Body: { qr_code: "...", new_status: "...", location_id: "..." }
- Looks up equipment by QR code
- Validates status transition
- Updates status
- This is the mobile-friendly quick-update endpoint for staff scanning QR codes

TASK 6 — Tests:
- WebSocket hub: test client registration, unregistration, broadcast to correct tenant only, no cross-tenant leak
- Redis bridge: mock Redis, test message routing to hub
- Availability state: test cache update on status change, test broadcast is triggered
- Board API: test cache hit vs DB fallback
```

---

## PHASE 3 — Sharing Request Workflow & Notifications

### Prompt:

```
Using the global context above, implement Phase 3 of MediFlow.
Phases 1 and 2 complete. Equipment catalogue, WebSocket hub, live board all work.

TASK 1 — Sharing Request Migrations:

Migration 009 — sharing_requests:
CREATE TABLE sharing_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  requesting_dept_id UUID NOT NULL REFERENCES departments(id),
  requesting_user_id UUID NOT NULL REFERENCES users(id),
  source_dept_id UUID REFERENCES departments(id),
  equipment_id UUID REFERENCES equipment_items(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  quantity_needed INTEGER NOT NULL DEFAULT 1,
  urgency VARCHAR(20) NOT NULL DEFAULT 'normal', -- low, normal, high, emergency
  reason TEXT,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  -- status: pending, matched, approved, declined, in_transit, active, return_pending, completed, cancelled
  needed_by TIMESTAMPTZ,
  expected_return_at TIMESTAMPTZ,
  matched_at TIMESTAMPTZ,
  approved_at TIMESTAMPTZ,
  approved_by_user_id UUID REFERENCES users(id),
  declined_reason TEXT,
  handoff_confirmed_by_source BOOLEAN DEFAULT false,
  handoff_confirmed_by_requester BOOLEAN DEFAULT false,
  handed_off_at TIMESTAMPTZ,
  return_confirmed_by_source BOOLEAN DEFAULT false,
  return_confirmed_by_requester BOOLEAN DEFAULT false,
  returned_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_requests_tenant ON sharing_requests(tenant_id);
CREATE INDEX idx_requests_status ON sharing_requests(status);
CREATE INDEX idx_requests_requesting_dept ON sharing_requests(requesting_dept_id);
CREATE INDEX idx_requests_source_dept ON sharing_requests(source_dept_id);

Migration 010 — request_history:
CREATE TABLE request_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  request_id UUID NOT NULL REFERENCES sharing_requests(id),
  changed_by_user_id UUID REFERENCES users(id),
  old_status VARCHAR(50),
  new_status VARCHAR(50),
  notes TEXT,
  changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 011 — notifications:
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  user_id UUID REFERENCES users(id), -- null = broadcast to whole department
  department_id UUID REFERENCES departments(id),
  type VARCHAR(100) NOT NULL,
  title VARCHAR(255) NOT NULL,
  message TEXT NOT NULL,
  data_json JSONB,
  is_read BOOLEAN NOT NULL DEFAULT false,
  read_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifications_user ON notifications(user_id, is_read);
CREATE INDEX idx_notifications_dept ON notifications(department_id, is_read);

TASK 2 — Sharing Request Service:
Implement backend/internal/request/ with full pattern:

Smart Matching Algorithm (MatchBestEquipment):
Given: requesting_dept_id, category_id, quantity_needed

Step 1: Find all available equipment of the requested category:
  - Status = available
  - is_shared = true
  - department_id != requesting_dept_id
  - Order candidates by score (higher = better match)

Step 2: Score each candidate:
  - Same building: +30 points
  - Same floor: +20 additional points
  - Department with excess stock (available > min_stock + 2): +25 points
  - Department that has received from requester before (reciprocal): +15 points
  - Lower utilisation rate in past 7 days: +10 points

Step 3: Return top match with source_dept and specific equipment_item

Service methods:
- CreateRequest(req): 
  1. Validate fields
  2. Run MatchBestEquipment
  3. If match found: set source_dept_id, equipment_id, status=matched, notify source dept
  4. If no match: status=pending (will be matched when equipment becomes available), notify requester
  5. Log history
  6. Publish WebSocket event

- ApproveRequest(requestID, approverUserID):
  1. Validate approver is from source department AND has role dept_head or charge_nurse
  2. Validate request status is matched
  3. Reserve the equipment item (status = reserved)
  4. Update request status = approved
  5. Notify requesting department
  6. Log history, publish event

- DeclineRequest(requestID, approverUserID, reason):
  1. Validate approver is from source department
  2. Update request status = declined
  3. Try to find next best match and re-assign
  4. Notify requesting department
  5. Log history

- ConfirmHandoff(requestID, confirmingUserID):
  1. Identify which side the user is on (source or requesting)
  2. Update appropriate handoff_confirmed field
  3. If both confirmed: update equipment status to in_use, department to requesting dept, status = active
  4. Notify both parties
  5. Publish WebSocket availability update

- ConfirmReturn(requestID, confirmingUserID):
  1. Same dual-confirm pattern as handoff
  2. If both confirmed: equipment returns to available in source department
  3. Update request status = completed
  4. Trigger utilisation log entry
  5. Publish WebSocket event

- CancelRequest(requestID, cancellingUserID):
  1. Only requester or admin can cancel
  2. If equipment was reserved: release reservation (status back to available)
  3. Update status = cancelled
  4. Notify source department

- GetPendingApprovalsForDepartment(deptID): requests awaiting approval from this dept
- GetActiveRequestsForDepartment(deptID): all non-terminal requests for this dept
- GetRequestHistory(filters, pagination)

API Endpoints:
- POST /api/v1/requests
- GET /api/v1/requests (paginated, filterable)
- GET /api/v1/requests/:id
- POST /api/v1/requests/:id/approve
- POST /api/v1/requests/:id/decline
- POST /api/v1/requests/:id/confirm-handoff
- POST /api/v1/requests/:id/confirm-return
- DELETE /api/v1/requests/:id (cancel)
- GET /api/v1/requests/pending-approvals (for current user's department)
- GET /api/v1/requests/active (active requests involving current user's department)

TASK 3 — Auto-Match Background Job:
Goroutine that runs every 5 minutes:
1. Find all requests with status=pending (no match yet)
2. Try to match each one again (equipment may have become available since)
3. If match found: update to matched, notify source department
4. Log job run with count of matches made

TASK 4 — Minimum Stock Alert System:
Background goroutine runs every 15 minutes:
1. For each tenant, for each department+category with a minimum_stock setting:
2. Count available equipment
3. If count < minimum: check if alert already sent in last 4 hours (Redis key: min_stock_alert:{dept_id}:{cat_id})
4. If not: create notification, publish to WebSocket, set Redis key with 4h TTL

Also check: when equipment becomes unavailable (via status update), immediately trigger a min-stock check for that department+category combination.

TASK 5 — Notification Service:
Implement backend/internal/alert/ with:

NotificationService methods:
- CreateNotification(tenantID, userID or deptID, type, title, message, data)
  - Saves to DB
  - Publishes to Redis: ws-events:{tenantID} with type="notification"
  - WebSocket bridge delivers to connected clients

- GetMyNotifications(userID, pagination): user's own + their department's unread notifications
- MarkRead(notificationID, userID)
- MarkAllRead(userID)
- GetUnreadCount(userID)

API Endpoints:
- GET /api/v1/notifications
- PUT /api/v1/notifications/:id/read
- PUT /api/v1/notifications/read-all
- GET /api/v1/notifications/count

SSE Endpoint (alternative to WebSocket for notifications only):
- GET /api/v1/notifications/stream
- Server-Sent Events stream
- On connection: send current unread count
- Push new notifications as they arrive

Also implement transit timeout checker:
Background goroutine every 30 minutes:
- Find equipment items with status=in_transit for more than 2 hours
- Find sharing requests with status=approved but handoff not confirmed in 3 hours
- Create warning notification to both departments
- Create "potentially lost" alert if over 24 hours

TASK 6 — Tests:
- Smart matching algorithm: test scoring with known data, verify best match selected
- Request state machine: all valid/invalid transitions
- Dual-confirm handoff: confirm partial, confirm both, out-of-order confirmations
- Auto-match job: seeds pending request, makes equipment available, verify match found
- Min stock alert: rate limiting (not re-alerted within 4h)
- Transit timeout: test alert fires at correct time thresholds
```

---

## PHASE 4 — Utilisation Tracking & Analytics

### Prompt:

```
Using the global context above, implement Phase 4 of MediFlow.
Phases 1-3 complete. Equipment, WebSocket, requests, notifications all work.

TASK 1 — Utilisation Tracking Migrations:

Migration 012 — utilisation_logs:
CREATE TABLE utilisation_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  equipment_id UUID NOT NULL REFERENCES equipment_items(id),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID REFERENCES departments(id),
  location_id UUID REFERENCES locations(id),
  status VARCHAR(50) NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  ended_at TIMESTAMPTZ,
  duration_minutes INTEGER, -- populated when ended_at is set
  sharing_request_id UUID REFERENCES sharing_requests(id), -- if usage was via sharing
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_utilisation_equipment ON utilisation_logs(equipment_id, started_at);
CREATE INDEX idx_utilisation_tenant ON utilisation_logs(tenant_id, started_at);
CREATE INDEX idx_utilisation_department ON utilisation_logs(department_id, started_at);

Migration 013 — demand_forecasts:
CREATE TABLE demand_forecasts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  forecast_date DATE NOT NULL,
  hour_of_day INTEGER, -- 0-23, null = daily forecast
  predicted_demand DECIMAL(5,2),
  confidence DECIMAL(5,2),
  model_version VARCHAR(50),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

Migration 014 — analytics_snapshots:
CREATE TABLE utilisation_snapshots (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id),
  snapshot_date DATE NOT NULL,
  department_id UUID NOT NULL REFERENCES departments(id),
  category_id UUID NOT NULL REFERENCES equipment_categories(id),
  total_items INTEGER,
  avg_utilisation_pct DECIMAL(5,2),
  total_in_use_minutes INTEGER,
  total_idle_minutes INTEGER,
  sharing_requests_received INTEGER,
  sharing_requests_sent INTEGER,
  UNIQUE(tenant_id, snapshot_date, department_id, category_id)
);

TASK 2 — Utilisation Tracking Service:
Implement backend/internal/utilisation/:

The utilisation_logs table is automatically maintained by the equipment status change flow:
- When status changes TO in_use: create a utilisation_log record with started_at=now, ended_at=null
- When status changes FROM in_use: find the open log (ended_at IS NULL), set ended_at=now, calculate duration_minutes

UtilisationService methods:
- GetEquipmentUtilisationRate(equipmentID, startDate, endDate):
  - Sum of duration_minutes where status=in_use in period / total minutes in period
  - Returns percentage

- GetDepartmentUtilisationSummary(tenantID, deptID, dateRange):
  - Per category: total items, avg utilisation rate, peak hours
  - Items sorted by utilisation rate (lowest first = idle candidates)

- GetHourlyPattern(tenantID, deptID, categoryID, lookbackDays):
  - For each hour 0-23, average utilisation rate across the lookback period
  - Returns array of 24 values — shows time-of-day patterns

- GetIdleEquipment(tenantID, minIdleDays):
  - Equipment that has not been in in_use status for at least minIdleDays
  - Grouped by department, sorted by idle duration descending

- GetTopSharedItems(tenantID, days):
  - Equipment items involved in most sharing requests
  - Indicates high-demand items that may warrant additional procurement

TASK 3 — Demand Forecasting:
Implement backend/internal/analytics/forecasting.go:

ForecastDemand(tenantID, deptID, categoryID, forecastDays):
This is a time-series forecasting using weighted moving average.

Step 1: Collect historical data
- For each day in last 90 days, count: number of in_use events, peak concurrent in_use count, sharing requests received

Step 2: Calculate day-of-week patterns
- Average demand for each day of week (Mon-Sun) over last 12 weeks
- day_of_week_factor[dow] = avg_demand_on_this_dow / overall_avg_demand

Step 3: Weighted Moving Average (WMA)
- Use last 30 days of daily demand
- Weights: most recent day = 30, next = 29, ... oldest = 1
- WMA = sum(demand[i] * weight[i]) / sum(weights)

Step 4: Apply day-of-week factor for each forecast day
- predicted_demand[day] = WMA * day_of_week_factor[day_of_week(forecast_date)]

Step 5: Confidence calculation
- Higher confidence if: more data, lower variance in recent history
- confidence = 1 - (std_dev / mean) — capped between 0.3 and 0.95

Store forecasts in demand_forecasts table.
Schedule: run forecast for all tenant+dept+category combinations weekly (every Monday 03:00 UTC).

API:
- GET /api/v1/analytics/demand-forecast?dept_id=&category_id=&days=7

TASK 4 — Analytics Service:
Implement backend/internal/analytics/ with all analytics endpoints:

OverviewDashboard(tenantID):
Returns:
- Total equipment count by status (available/in_use/maintenance/etc)
- Overall utilisation rate (today, this week, this month trend)
- Active sharing requests count
- Pending approvals count
- Top 5 idle equipment
- Top 5 overutilised categories (where utilisation > 85%)
- Recent alerts (last 10)

UtilisationReport(tenantID, dateRange, groupBy: department|category|item):
- Utilisation rates grouped by specified dimension
- Trend chart data: daily utilisation rates over the period
- Returns data formatted for recharts bar/line charts

SharingReport(tenantID, dateRange):
- Total sharing requests: by status, by urgency
- Average approval time
- Average handoff time (approval → handoff confirmed)
- Average return time (active → completed)
- Department-level sharing activity matrix: how much each dept gives to and receives from each other
- Top requested categories

IdlenessReport(tenantID):
- Equipment sorted by idle days descending
- Estimated cost of idle equipment (purchase_cost * idle_days / expected_lifecycle_days)
- Redeployment suggestions: idle equipment in dept A, high demand in dept B

ProcurementInsightsReport(tenantID):
- Categories with > 85% avg utilisation: recommend purchasing more
- Categories with < 20% avg utilisation across all departments: recommend reviewing inventory
- Equipment items nearing end of life (age > 80% of expected lifecycle): replacement candidates

DepartmentComparison(tenantID):
- For each department: avg utilisation, shares given, shares received, net sharing (donor vs receiver)
- Visual-friendly data for heatmap

API Endpoints (all cached in Redis, TTL varies):
- GET /api/v1/analytics/overview (TTL: 5 min)
- GET /api/v1/analytics/utilisation?start=&end=&group_by=
- GET /api/v1/analytics/sharing?start=&end=
- GET /api/v1/analytics/idleness
- GET /api/v1/analytics/procurement-insights
- GET /api/v1/analytics/department-comparison
- GET /api/v1/analytics/demand-forecast?dept_id=&category_id=&days=
- GET /api/v1/analytics/hourly-pattern?dept_id=&category_id=

Export:
- GET /api/v1/reports/equipment-list (CSV)
- GET /api/v1/reports/utilisation (CSV)
- GET /api/v1/reports/sharing-history (CSV)

TASK 5 — Daily Snapshot Job:
Goroutine at midnight UTC:
- For each tenant, dept, category combination with equipment:
  - Calculate today's utilisation metrics
  - Upsert into utilisation_snapshots
  - Cache is cleared for that tenant

TASK 6 — Tests:
- Utilisation rate calculation: known start/end times, verify exact percentage
- Hourly pattern: seed data, verify pattern matches expected
- WMA forecasting: seed known time series, verify forecast is correct
- Idle equipment: seed equipment with no recent usage, verify correct idle days
- Procurement insights: threshold conditions properly identified
```

---

## PHASE 5 — NextJS Frontend

### Prompt:

```
Using the global context above, implement Phase 5 of MediFlow — the complete NextJS frontend.

Design direction: Operational, real-time feel. Clean dashboard like hospital operations centre.
- Colours: Primary #1D4ED8 (blue), Available #16A34A (green), InUse #D97706 (amber), Maintenance #6B7280 (grey), Critical/Missing #DC2626 (red), Background #F1F5F9
- Font: Inter
- Use shadcn/ui throughout
- Responsive for desktop and tablet

TASK 1 — WebSocket Client Setup:
Create lib/websocket.ts:
- WebSocket client singleton using native browser WebSocket
- Auto-reconnect with exponential backoff (max 30s between retries)
- Message dispatcher: on message, parse JSON and dispatch to registered handlers by type
- Zustand store: wsStore with connected state, lastMessage, and subscribe/unsubscribe for specific message types
- On availability_update message: update availability board store
- On notification message: add to notifications store, increment unread count

TASK 2 — App Structure:
app/
├── (auth)/login/page.tsx
├── (dashboard)/
│   ├── layout.tsx
│   ├── page.tsx (overview)
│   ├── board/page.tsx (live availability board)
│   ├── equipment/
│   │   ├── page.tsx
│   │   ├── [id]/page.tsx
│   │   └── new/page.tsx
│   ├── requests/
│   │   ├── page.tsx
│   │   ├── new/page.tsx
│   │   └── [id]/page.tsx
│   ├── analytics/page.tsx
│   └── settings/
│       ├── departments/page.tsx
│       ├── categories/page.tsx
│       └── users/page.tsx

TASK 3 — Live Availability Board (Most Important Page):
This is the flagship feature. Make it visually impressive and functional.

Layout: Full-width table/grid
- Column headers = equipment categories (Ventilator, ECG Machine, Infusion Pump, etc.) with icons
- Row headers = departments (ICU, Post-Op, Emergency, General Ward, etc.)
- Each cell = availability summary for that dept × category combination

Each cell shows:
- Large number: available count (green if > min_stock, amber if == min_stock, red if 0)
- Small text: "3 of 8 available"
- Small text: "2 in use"
- Click to expand: shows individual device names and statuses

Real-time updates:
- WebSocket message of type availability_update → find the matching cell → animate it updating
- Use CSS transition so the number change is visible (brief highlight flash)
- Small "live" indicator badge in top right of page (green dot pulsing if connected)

Filter bar above board:
- Filter by building, floor
- Toggle: show only categories with low availability
- Search category

Cell click modal/drawer:
- Shows all individual equipment items for that dept × category
- Each item: name, status badge, current location, last status change time
- "Request this equipment" button (appears if status=available)
- Engineer can click to quickly update status

TASK 4 — Dashboard Overview Page:
4 stat cards row:
- Total Equipment | Available Now | In Use | Needs Attention (maintenance + missing)

Live activity feed (right side of page):
- Last 20 equipment status changes (real-time via WebSocket)
- Each entry: device name, old status → new status, department, time ago
- New entries slide in from top

Left side:
- Sharing requests summary: pending approvals (with action button), active transfers
- Equipment needing attention: in_maintenance or missing items
- Departments with low stock alerts

TASK 5 — Sharing Requests Pages:
List page:
- Three tabs: My Department's Requests | Pending My Approval | All Requests
- Each request card: equipment category, requesting dept → source dept, urgency badge, status badge, time, actions
- Status badges: pending (grey), matched (blue), approved (teal), in_transit (orange), active (green), completed (grey-green), declined (red)
- Action buttons context-aware: Approve/Decline for source dept head, Confirm Handoff, Confirm Return

New Request page:
- Step 1: Select category + quantity + urgency + reason + needed by date
- Step 2: System shows matched source department + specific equipment (from API response)
- Step 3: Confirm submission
- Show "no match found — request will be queued" if no match

Request Detail page:
- Full request timeline (history of all status changes)
- Equipment details card
- Both-department confirmation UI for handoff and return (show who has confirmed, who hasn't)
- Decline reason if declined

TASK 6 — Equipment Pages:
List page:
- Table with: name, category icon, department, current location, status badge (colour coded), last update time
- Filter by status, category, department
- Click status badge → inline quick-update dropdown (if authorised)
- QR code icon → show QR code in modal

Detail page:
- Equipment info
- Current status with large coloured badge
- Current location map placeholder (show floor/room as text)
- Status history timeline (full vertical timeline component)
- Active sharing request for this device (if any)
- Quick status update panel (role-dependent)

TASK 7 — Analytics Page:
Tab navigation: Overview | Utilisation | Sharing | Idleness | Procurement

Overview tab:
- Utilisation trend line chart (this week vs last week)
- Department utilisation heatmap: table where cells are colour-coded by utilisation %
  - rows = departments, columns = categories, cell colour from green (high) to red (low)
- Active requests funnel: pending → approved → in_transit → active → completed

Utilisation tab:
- Date range picker
- Bar chart: utilisation by department
- Line chart: utilisation trend over time
- Hourly pattern chart: 24-bar chart showing when equipment is most used

Sharing tab:
- Sankey diagram placeholder (or simple table): which departments share most with which
- Avg approval/handoff/return time cards
- Most requested categories bar chart

Idleness tab:
- Table: idle equipment sorted by idle days, with estimated idle cost
- "Redeploy" action: creates a suggestion notification to the department (stub)

Procurement tab:
- Two sections: "Consider Buying More" (high utilisation) and "Excess Inventory" (low utilisation)
- Each with category name, avg utilisation %, recommendation

TASK 8 — Notifications UI:
Bell icon in topnav with unread count badge.
Click → slide-in drawer showing notifications:
- Each notification: type icon, title, message, time ago, read/unread indicator
- Click to mark read
- "Mark all read" button
- Link to relevant page if applicable (request ID, equipment ID)

Real-time: new notifications from WebSocket push immediately appear at top of list, badge count increments.

TASK 9 — React Query & API Layer:
lib/api/:
- useAvailabilityBoard, useEquipmentDetail, useEquipmentList
- useCreateRequest, useApproveRequest, useDeclineRequest, useConfirmHandoff, useConfirmReturn
- useRequestList, useRequestDetail, usePendingApprovals
- useNotifications, useUnreadCount, useMarkRead
- useAnalyticsOverview, useUtilisationReport, useSharingReport, useIdlenessReport

Zustand stores:
- boardStore: availability board state, updated by WebSocket messages
- notificationStore: notifications list, unread count, updated by WebSocket
- wsStore: connection status, reconnecting state
```

---

## PHASE 6 — DevOps, Load Testing & Documentation

### Prompt:

```
Using the global context above, implement Phase 6 of MediFlow.

TASK 1 — Multi-stage Dockerfiles:
Backend: golang:1.22-alpine builder → alpine:3.19 runner with non-root user
Frontend: node:20-alpine deps → builder → runner (Next standalone output)

TASK 2 — Kubernetes Manifests (k8s/):
- namespace.yaml
- postgres/: StatefulSet, PVC, Service, Secret
- redis/: Deployment (with keyspace notification config), Service
- backend/: Deployment (3 replicas), Service, HPA (CPU 70% → max 10), ConfigMap, Secret
- frontend/: Deployment (2 replicas), Service
- ingress.yaml: route /api/* and /ws to backend, /* to frontend, WebSocket upgrade headers configured
- Note in README: WebSocket requires nginx ingress with proxy_read_timeout 3600 and proxy_send_timeout 3600

TASK 3 — GitHub Actions:
.github/workflows/ci.yml:
- test-backend: postgres+redis services, go test -race ./...
- test-frontend: lint + build
- build-and-push (main branch only): build + push both images to GHCR

TASK 4 — WebSocket Load Test (k6):
k6/ws_load_test.js:
- 200 virtual users each open a WebSocket connection
- Hold connection for 60 seconds
- Every 10 seconds: simulate an equipment status update via REST API (10 users do this concurrently)
- Measure: connection establishment time, message delivery latency, reconnection success rate

k6/rest_load_test.js:
- 300 virtual users
- Mix of: GET /availability-board (40%), GET /requests (30%), POST /requests (20%), PUT /status (10%)
- Thresholds: p95 < 150ms, error rate < 0.5%

Document results (paste actual numbers) in README.

TASK 5 — Seed Data:
backend/cmd/seed/main.go:
- 1 demo tenant: "City General Hospital"
- 5 departments: ICU, Emergency, Post-Op, General Ward, Cardiology
- 8 equipment categories: Ventilator, ECG Machine, Infusion Pump, Portable Ultrasound, Defibrillator, Cardiac Monitor, Pulse Oximeter, Wheelchair
- 3-8 equipment items per category per department (realistic quantities)
- Set minimum stock levels for each dept × category
- 6 months of realistic utilisation history (vary by time of day and day of week)
- 50 historical sharing requests with realistic statuses
- All health scores and snapshots pre-calculated
- Demo users: one per role

TASK 6 — README:
Sections:
1. Banner with project name and tagline: "Real-time equipment sharing and utilisation intelligence for hospitals"
2. Problem statement
3. Live demo link (if deployed)
4. Feature highlights with key screenshots described
5. Architecture diagram (ASCII)
6. Tech stack table
7. Getting started: prerequisites + `docker-compose up --build` one-liner
8. WebSocket architecture explanation (brief)
9. Kubernetes deployment with WebSocket note
10. Load test results
11. Project structure
12. Future roadmap: IoT integration, ML forecasting, mobile app
13. License

TASK 7 — Final Checklist:
- [ ] WebSocket reconnection works after backend restart
- [ ] Tenant isolation verified: user from tenant A cannot see tenant B's board
- [ ] All status transitions enforced
- [ ] Dual-confirm handoff and return work correctly
- [ ] Min stock alerts rate-limited correctly
- [ ] Analytics endpoints return chart-ready data
- [ ] Seed data produces a realistic populated dashboard
- [ ] Docker Compose starts clean with docker-compose up
- [ ] Swagger accessible and complete
- [ ] CI pipeline passes
```
