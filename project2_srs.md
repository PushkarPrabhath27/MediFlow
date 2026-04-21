# Software Requirements Specification
## MediFlow – Hospital Equipment Utilisation & Sharing Optimiser

**Document Version:** 1.0  
**Date:** April 2025  
**Status:** Approved  
**Author:** [Your Name]

---

## Table of Contents
1. Introduction
2. Overall Description
3. Stakeholders
4. Functional Requirements
5. Non-Functional Requirements
6. System Architecture Requirements
7. Data Requirements
8. API Requirements
9. Security Requirements
10. Integration Requirements
11. Constraints & Assumptions
12. Glossary

---

## 1. Introduction

### 1.1 Purpose
This document specifies the software requirements for MediFlow, a real-time multi-tenant SaaS platform that enables hospitals to track equipment availability, optimise utilisation, and manage inter-department equipment sharing. It serves as the authoritative reference for design, development, and testing.

### 1.2 Scope
MediFlow encompasses:
- A web-based real-time availability board for all shared equipment
- A structured equipment sharing request and approval workflow
- A WebSocket-based real-time update delivery system
- A utilisation tracking and analytics engine
- A demand forecasting module
- A notification and alert system
- Multi-tenancy for multiple hospital organisations

Out of scope for v1.0:
- Mobile native applications (web is responsive)
- RFID or IoT-based automatic location tracking
- Integration with existing Hospital Information Systems (HIS)
- Patient data of any kind

### 1.3 Definitions

| Term | Definition |
|---|---|
| Tenant | A hospital or healthcare facility using MediFlow |
| Equipment Item | A single physical device tracked in the system |
| Availability Board | Real-time grid showing equipment availability per dept × category |
| Sharing Request | A formal digital request by one department to borrow equipment from another |
| Utilisation Rate | % of time a device is actively in use over a given period |
| Dual-Confirm | A handoff or return process that requires both parties to confirm |
| Min Stock Level | The minimum number of available items a department wants to maintain |
| WebSocket Hub | The GoLang service managing all persistent WebSocket connections |

### 1.4 References
- Affordmed Full Stack Developer Job Description
- WHO Technical Series on Safe Management of Wastes from Health-Care Activities
- ECRI Institute Equipment Management Handbook

---

## 2. Overall Description

### 2.1 Product Perspective
MediFlow is a standalone SaaS product with no required integration with existing hospital IT systems in v1.0. Its REST API is designed for future integration capability.

### 2.2 Product Functions Summary
- Real-time equipment availability tracking across departments
- WebSocket-based live board updates
- Inter-department equipment sharing request workflow
- Automated smart matching of requests to available equipment
- Utilisation logging and analytics
- Demand forecasting
- Minimum stock level alerts
- Role-based access control
- Multi-tenancy

### 2.3 Operating Environment
- Backend deployed on Kubernetes (Azure AKS or equivalent)
- Frontend deployed as containerised NextJS application (standalone output mode)
- PostgreSQL 15 as primary data store
- Redis 7 with keyspace notification enabled (for pub/sub)
- Nginx ingress configured with WebSocket upgrade headers
- Accessible via modern web browsers (Chrome 90+, Firefox 88+, Safari 14+, Edge 90+)

---

## 3. Stakeholders

| Stakeholder | Role | Primary Concerns |
|---|---|---|
| Charge Nurse / Ward Sister | Daily user — manages ward equipment | Check availability, create requests, confirm handoffs |
| Department Head | Departmental authority | Approve/decline requests, see utilisation for their dept |
| Biomedical Engineer | Maintains equipment | Track locations, update maintenance status |
| Hospital Administrator | Facility-wide oversight | Cross-department view, procurement insights |
| Procurement Officer | Equipment purchasing | Utilisation data to guide procurement decisions |

---

## 4. Functional Requirements

### 4.1 Authentication & Authorisation

**FR-AUTH-001:** The system shall support email and password based login with JWT tokens.  
**FR-AUTH-002:** Access tokens shall expire after 15 minutes; refresh tokens after 7 days.  
**FR-AUTH-003:** The system shall enforce role-based access control with six roles: super_admin, hospital_admin, department_head, charge_nurse, staff, engineer.  
**FR-AUTH-004:** Role permission matrix:

| Action | staff | engineer | charge_nurse | dept_head | hosp_admin | super_admin |
|---|---|---|---|---|---|---|
| View availability board | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Update equipment status | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Create sharing request | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Approve sharing request | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Confirm handoff/return | ❌ | ❌ | ✅ | ✅ | ✅ | ✅ |
| Manage equipment | ❌ | ✅ | ❌ | ✅ | ✅ | ✅ |
| View analytics | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Manage users/departments | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ |

### 4.2 Tenant Management

**FR-TEN-001:** The system shall support multiple fully isolated tenants.  
**FR-TEN-002:** Tenant plans: basic (≤100 devices, ≤10 departments), pro (≤1000 devices, ≤50 departments), enterprise (unlimited).  
**FR-TEN-003:** Tenant data shall be completely isolated.  
**FR-TEN-004:** Super admins shall manage tenant lifecycle.

### 4.3 Equipment Catalogue

**FR-EQ-001:** Authorised users shall be able to register equipment with: name (required), category (required), department, location, model, manufacturer, serial number, asset tag, purchase date, purchase cost, sharing flag.  
**FR-EQ-002:** The system shall auto-generate a unique QR code for each equipment item.  
**FR-EQ-003:** Equipment shall have a status field with the following valid values: available, in_use, reserved, in_maintenance, in_transit, missing, decommissioned.  
**FR-EQ-004:** Equipment shall track its current department and location.  
**FR-EQ-005:** Equipment with is_shared = false shall not appear in sharing request matching.  
**FR-EQ-006:** The system shall support bulk CSV import for equipment items.  
**FR-EQ-007:** Full CRUD shall be available for equipment categories and locations.

### 4.4 Equipment Status Management

**FR-ST-001:** Status transitions shall follow a defined state machine. Valid transitions:

| From | Valid Next States |
|---|---|
| available | in_use, reserved, in_maintenance, in_transit, missing |
| in_use | available, in_maintenance, missing |
| reserved | in_use, available |
| in_maintenance | available, decommissioned |
| in_transit | available, in_use |
| missing | available, decommissioned |

**FR-ST-002:** Every status change shall be recorded in equipment_status_logs with: old status, new status, changed by user, reason, timestamp.  
**FR-ST-003:** A QR scan endpoint shall allow quick status updates by scanning a device's QR code.  
**FR-ST-004:** Status updates shall trigger: Redis availability cache update, WebSocket broadcast to tenant clients.

### 4.5 Live Availability Board

**FR-BOARD-001:** The system shall provide a real-time availability board showing equipment availability at the intersection of department × category.  
**FR-BOARD-002:** Each board cell shall display: available count, in-use count, total count, below-minimum indicator.  
**FR-BOARD-003:** Board data shall be served from Redis cache (TTL 90 seconds, refreshed every 60 seconds).  
**FR-BOARD-004:** When equipment status changes, the affected board cell shall be updated in Redis and a WebSocket event published within 500ms.  
**FR-BOARD-005:** All connected clients for the tenant shall receive the board update within 1 second of the status change.  
**FR-BOARD-006:** The board API endpoint shall fall back to direct DB query if Redis is unavailable.

### 4.6 WebSocket System

**FR-WS-001:** The system shall maintain persistent WebSocket connections with authenticated clients.  
**FR-WS-002:** JWT authentication shall be required to establish a WebSocket connection, passed as a URL query parameter.  
**FR-WS-003:** WebSocket messages shall be tenant-isolated: a client shall only receive events for their own tenant.  
**FR-WS-004:** The system shall implement ping/pong keepalive on WebSocket connections (30-second ping interval).  
**FR-WS-005:** The system shall support automatic client reconnection with exponential backoff (up to 30 seconds).  
**FR-WS-006:** WebSocket message types: availability_update, notification, request_update.  
**FR-WS-007:** All WebSocket events shall be routed through Redis pub/sub to support multiple backend instances.

### 4.7 Equipment Sharing Request Workflow

**FR-REQ-001:** Authorised users shall be able to create sharing requests specifying: category, quantity, urgency (low/normal/high/emergency), reason, needed-by datetime, expected return datetime.  
**FR-REQ-002:** Upon request creation, the system shall automatically run a smart matching algorithm to identify the best available source department and equipment item.  
**FR-REQ-003:** Smart matching scoring criteria (higher score = better match):
- Same building as requesting department: +30 points
- Same floor: +20 points
- Department has excess stock (available > min_stock + 2): +25 points
- Reciprocal historical relationship: +15 points
- Lower recent utilisation rate: +10 points

**FR-REQ-004:** If a match is found, request status becomes matched and the source department is notified.  
**FR-REQ-005:** If no match found, request status becomes pending and a background job will retry matching every 5 minutes.  
**FR-REQ-006:** The source department (role: dept_head or charge_nurse) shall be able to approve or decline a matched request.  
**FR-REQ-007:** On approval, the matched equipment item status changes to reserved.  
**FR-REQ-008:** When equipment is physically transferred, both source and requesting departments must individually confirm the handoff.  
**FR-REQ-009:** Upon both parties confirming handoff: equipment status changes to in_use, department ownership changes to requesting department, request status changes to active.  
**FR-REQ-010:** When equipment is returned, both departments must individually confirm the return.  
**FR-REQ-011:** Upon both parties confirming return: equipment status returns to available in source department, request status changes to completed, utilisation log entry is closed.  
**FR-REQ-012:** Request status state machine:

| From | Valid Next States |
|---|---|
| pending | matched, cancelled |
| matched | approved, declined, cancelled |
| approved | in_transit (on handoff start), cancelled |
| in_transit | active (on handoff complete) |
| active | return_pending |
| return_pending | completed (on return complete) |
| declined | (terminal) |
| completed | (terminal) |
| cancelled | (terminal) |

**FR-REQ-013:** Every status change shall be logged in request_history.  
**FR-REQ-014:** A background job shall check for requests with status=pending every 5 minutes and attempt re-matching.

### 4.8 Minimum Stock Level Alerts

**FR-MSL-001:** Department heads shall be able to set minimum stock levels per equipment category for their department.  
**FR-MSL-002:** The system shall check stock levels every 15 minutes and when any status change makes an item unavailable.  
**FR-MSL-003:** When available count drops below minimum, an alert shall be generated with severity=warning.  
**FR-MSL-004:** Minimum stock alerts shall be rate-limited to one alert per department+category combination per 4-hour window.  
**FR-MSL-005:** The alert shall include the system's suggestion for which department has surplus equipment that could be requested.

### 4.9 Transit Timeout Monitoring

**FR-TT-001:** The system shall monitor equipment items with status=in_transit.  
**FR-TT-002:** If an item remains in_transit for more than 2 hours, a warning notification shall be sent to both involved departments.  
**FR-TT-003:** If an item remains in_transit for more than 24 hours, a critical alert shall be raised flagging potential equipment loss.  
**FR-TT-004:** Transit timeout monitoring shall run every 30 minutes.

### 4.10 Notification System

**FR-NOT-001:** The system shall create notifications for the following events: request created, request matched, request approved, request declined, handoff confirmed (both sides), return confirmed (both sides), min stock alert, transit timeout warning, critical transit alert.  
**FR-NOT-002:** Notifications shall target: a specific user, all users in a department, or all hospital_admins.  
**FR-NOT-003:** Notifications shall be delivered in real-time via WebSocket to connected clients.  
**FR-NOT-004:** A Server-Sent Events (SSE) fallback endpoint shall be provided for notification delivery.  
**FR-NOT-005:** Users shall be able to view, mark as read, and mark all as read for their notifications.  
**FR-NOT-006:** Unread notification count shall be retrievable for badge display.

### 4.11 Utilisation Tracking

**FR-UT-001:** The system shall automatically create a utilisation log entry whenever equipment status changes to in_use, recording start time, department, and location.  
**FR-UT-002:** When status changes from in_use to any other status, the open utilisation log entry shall be closed with end time and duration calculated.  
**FR-UT-003:** Utilisation logs shall record whether the usage occurred via a sharing request.  
**FR-UT-004:** A daily snapshot job shall run at midnight UTC to aggregate utilisation metrics by tenant + department + category.

### 4.12 Analytics & Reporting

**FR-AN-001:** Overview dashboard shall show: equipment by status, overall utilisation rates (today/week/month), active requests, pending approvals, departments with low stock, recent activity.  
**FR-AN-002:** Utilisation report shall show: rates by department and category, trend over time, hourly usage patterns.  
**FR-AN-003:** Sharing report shall show: request volumes by status and urgency, avg approval/handoff/return times, department sharing network.  
**FR-AN-004:** Idleness report shall show: equipment sorted by idle days, estimated idle cost, redeployment suggestions.  
**FR-AN-005:** Procurement insights report shall identify: categories with >85% avg utilisation (recommend buying more), categories with <20% utilisation (excess inventory), end-of-life items.  
**FR-AN-006:** Demand forecast shall predict equipment needs for each dept × category for the next 7 days.  
**FR-AN-007:** All analytics data shall be exportable as CSV.

---

## 5. Non-Functional Requirements

### 5.1 Performance

**NFR-PERF-001:** REST API response time shall be ≤ 100ms at p95 under normal load.  
**NFR-PERF-002:** Availability board API (from Redis cache) shall respond in ≤ 30ms.  
**NFR-PERF-003:** WebSocket message delivery latency from status change event to client receipt shall be ≤ 1 second.  
**NFR-PERF-004:** The system shall support 500 concurrent WebSocket connections per backend pod.  
**NFR-PERF-005:** The system shall handle 300 concurrent REST API users without degradation.  
**NFR-PERF-006:** Analytics report generation shall complete in ≤ 2 seconds with Redis caching.

### 5.2 Scalability

**NFR-SCAL-001:** Backend shall be horizontally scalable via Kubernetes HPA.  
**NFR-SCAL-002:** All WebSocket events shall route through Redis pub/sub to support multiple backend instances correctly.  
**NFR-SCAL-003:** No in-process state shall be used for data that must survive a pod restart.  
**NFR-SCAL-004:** HPA target: scale when CPU > 70%, max 10 replicas.

### 5.3 Reliability

**NFR-REL-001:** WebSocket connections shall implement auto-reconnection with exponential backoff on the client side.  
**NFR-REL-002:** Redis pub/sub subscriber shall implement reconnection logic with retry every 5 seconds on failure.  
**NFR-REL-003:** All background jobs shall be idempotent.  
**NFR-REL-004:** The availability board shall degrade gracefully to DB queries if Redis is unavailable.  
**NFR-REL-005:** The system shall implement graceful shutdown, closing all WebSocket connections cleanly.

### 5.4 Security

**NFR-SEC-001:** All passwords hashed with bcrypt, cost ≥ 12.  
**NFR-SEC-002:** All HTTP communication shall use HTTPS (TLS 1.2+).  
**NFR-SEC-003:** WebSocket connections shall require valid JWT authentication.  
**NFR-SEC-004:** Tenant isolation shall be enforced at the database query level on every query.  
**NFR-SEC-005:** Rate limiting: 100 req/min per IP globally, 10 login attempts/min per IP.  
**NFR-SEC-006:** All queries use parameterised statements.  
**NFR-SEC-007:** No secrets committed to source control.  
**NFR-SEC-008:** CORS restricted to whitelisted frontend origin.  
**NFR-SEC-009:** Nginx WebSocket proxy shall set appropriate timeout headers to prevent abuse.

### 5.5 Maintainability

**NFR-MAINT-001:** Repository/service/handler layered architecture in GoLang.  
**NFR-MAINT-002:** Unit test coverage ≥ 70% for business logic.  
**NFR-MAINT-003:** All API endpoints fully documented in Swagger.  
**NFR-MAINT-004:** Database changes managed via versioned migrations.  
**NFR-MAINT-005:** Structured JSON logging with log levels.

### 5.6 Usability

**NFR-USE-001:** Availability board shall display data within 1 second of page load.  
**NFR-USE-002:** Real-time cell updates shall be visually indicated (highlight animation).  
**NFR-USE-003:** Connection status (live/reconnecting/offline) shall be clearly shown on the board.  
**NFR-USE-004:** All forms shall have inline validation.

---

## 6. System Architecture Requirements

**AR-001:** Backend in GoLang with chi router and gorilla/websocket.  
**AR-002:** Frontend in NextJS 14 with App Router.  
**AR-003:** PostgreSQL 15 primary data store.  
**AR-004:** Redis 7 with keyspace notifications enabled.  
**AR-005:** WebSocket hub implemented as a goroutine-based broadcaster with tenant-level isolation.  
**AR-006:** All WebSocket events flow through Redis pub/sub to enable multi-pod deployment.  
**AR-007:** All services containerised with multi-stage Docker builds.  
**AR-008:** Kubernetes manifests provided for all components.  
**AR-009:** Nginx ingress must be configured with WebSocket proxy headers.  
**AR-010:** CI/CD via GitHub Actions.

---

## 7. Data Requirements

### 7.1 Data Retention
- Equipment status logs: retain 2 years
- Sharing request history: retain indefinitely
- Utilisation logs: retain 2 years, archive older
- Utilisation snapshots: retain indefinitely (aggregated data, small volume)
- Notifications: retain 6 months

### 7.2 Data Integrity
**DR-001:** All tenant IDs validated on every query.  
**DR-002:** Foreign key constraints enforced at DB level.  
**DR-003:** Soft deletes for equipment and users.  
**DR-004:** All timestamps stored in UTC.  
**DR-005:** Utilisation log entries must always be closed before a new one is opened for the same device.

---

## 8. API Requirements

**API-001:** RESTful naming conventions throughout.  
**API-002:** Consistent JSON envelope:
```json
{
  "success": true,
  "data": {},
  "error": null,
  "request_id": "uuid",
  "timestamp": "ISO8601"
}
```
**API-003:** WebSocket messages use typed envelopes:
```json
{
  "type": "availability_update | notification | request_update",
  "tenant_id": "uuid",
  "payload": {},
  "timestamp": "ISO8601"
}
```
**API-004:** Pagination on all list endpoints.  
**API-005:** API versioning via /api/v1/ prefix.  
**API-006:** Swagger at /swagger/index.html.

---

## 9. Security Requirements

Same as section 5.4. Additional:  
**SR-001:** WebSocket JWT validation must check token expiry and tenant membership.  
**SR-002:** A user connecting via WebSocket shall only receive events for their own tenant.

---

## 10. Integration Requirements

**IR-001:** The REST API is designed to accept future IoT status update webhooks at POST /api/v1/equipment/iot-update (stub endpoint in v1.0 with documentation).  
**IR-002:** The notification system is abstracted behind an interface to support future email/SMS integration.  
**IR-003:** CSV import/export is the primary data exchange mechanism in v1.0.

---

## 11. Constraints & Assumptions

### Constraints
- v1.0 uses weighted moving average for demand forecasting, not a full ML model.
- No physical location tracking — location is manually updated by staff.
- WebSocket requires nginx ingress to be configured with correct timeout/upgrade headers.
- Redis keyspace notifications must be enabled with config `notify-keyspace-events KEA`.

### Assumptions
- Hospital staff have access to a desktop or tablet browser.
- Each hospital is a single tenant.
- Equipment location is updated manually when items move between departments.
- No patient data is stored or referenced at any point.

---

## 12. Glossary

| Term | Meaning |
|---|---|
| Availability Board | Real-time grid showing equipment by dept × category |
| WebSocket Hub | GoLang goroutine managing all live client connections |
| Redis Bridge | Goroutine that subscribes to Redis and routes messages to WebSocket hub |
| Dual-Confirm | Handoff/return process requiring confirmation from both parties |
| Smart Matching | Algorithm that scores and selects the best source department for a sharing request |
| Utilisation Rate | % of time a device is in active use over a given period |
| Transit Timeout | Alert raised when equipment stays in_transit beyond expected time |
| Min Stock Alert | Warning when a department's available count drops below their set minimum |
| WMA | Weighted Moving Average — used for demand forecasting |
| HPA | Horizontal Pod Autoscaler (Kubernetes) |
| SSE | Server-Sent Events — server-to-client HTTP streaming (fallback to WebSocket) |
