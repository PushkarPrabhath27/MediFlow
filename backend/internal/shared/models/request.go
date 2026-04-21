package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	RequestStatusPending       RequestStatus = "pending"
	RequestStatusMatched       RequestStatus = "matched"
	RequestStatusApproved      RequestStatus = "approved"
	RequestStatusDeclined      RequestStatus = "declined"
	RequestStatusInTransit     RequestStatus = "in_transit"
	RequestStatusActive        RequestStatus = "active"
	RequestStatusReturnPending RequestStatus = "return_pending"
	RequestStatusCompleted     RequestStatus = "completed"
	RequestStatusCancelled     RequestStatus = "cancelled"
)

type SharingRequest struct {
	ID                         uuid.UUID     `db:"id" json:"id"`
	TenantID                   uuid.UUID     `db:"tenant_id" json:"tenant_id"`
	RequestingDeptID           uuid.UUID     `db:"requesting_dept_id" json:"requesting_dept_id"`
	RequestingUserID           uuid.UUID     `db:"requesting_user_id" json:"requesting_user_id"`
	SourceDeptID               *uuid.UUID    `db:"source_dept_id" json:"source_dept_id"`
	EquipmentID                *uuid.UUID    `db:"equipment_id" json:"equipment_id"`
	CategoryID                 uuid.UUID     `db:"category_id" json:"category_id"`
	QuantityNeeded             int           `db:"quantity_needed" json:"quantity_needed"`
	Urgency                    string        `db:"urgency" json:"urgency"`
	Reason                     string        `db:"reason" json:"reason"`
	Status                     RequestStatus `db:"status" json:"status"`
	NeededBy                   *time.Time    `db:"needed_by" json:"needed_by"`
	ExpectedReturnAt           *time.Time    `db:"expected_return_at" json:"expected_return_at"`
	MatchedAt                  *time.Time    `db:"matched_at" json:"matched_at"`
	ApprovedAt                 *time.Time    `db:"approved_at" json:"approved_at"`
	ApprovedByUserID           *uuid.UUID    `db:"approved_by_user_id" json:"approved_by_user_id"`
	DeclinedReason             string        `db:"declined_reason" json:"declined_reason"`
	HandoffConfirmedBySource   bool          `db:"handoff_confirmed_by_source" json:"handoff_confirmed_by_source"`
	HandoffConfirmedByRequester bool          `db:"handoff_confirmed_by_requester" json:"handoff_confirmed_by_requester"`
	HandedOffAt                *time.Time    `db:"handed_off_at" json:"handed_off_at"`
	ReturnConfirmedBySource    bool          `db:"return_confirmed_by_source" json:"return_confirmed_by_source"`
	ReturnConfirmedByRequester bool          `db:"return_confirmed_by_requester" json:"return_confirmed_by_requester"`
	ReturnedAt                 *time.Time    `db:"returned_at" json:"returned_at"`
	CreatedAt                  time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt                  time.Time     `db:"updated_at" json:"updated_at"`
}

type RequestHistory struct {
	ID              uuid.UUID     `db:"id" json:"id"`
	RequestID       uuid.UUID     `db:"request_id" json:"request_id"`
	ChangedByUserID *uuid.UUID    `db:"changed_by_user_id" json:"changed_by_user_id"`
	OldStatus       RequestStatus `db:"old_status" json:"old_status"`
	NewStatus       RequestStatus `db:"new_status" json:"new_status"`
	Notes           string        `db:"notes" json:"notes"`
	ChangedAt       time.Time     `db:"changed_at" json:"changed_at"`
}

type Notification struct {
	ID           uuid.UUID `db:"id" json:"id"`
	TenantID     uuid.UUID `db:"tenant_id" json:"tenant_id"`
	UserID       *uuid.UUID `db:"user_id" json:"user_id"`
	DepartmentID *uuid.UUID `db:"department_id" json:"department_id"`
	Type         string    `db:"type" json:"type"`
	Title        string    `db:"title" json:"title"`
	Message      string    `db:"message" json:"message"`
	DataJSON     []byte    `db:"data_json" json:"data_json"`
	IsRead       bool      `db:"is_read" json:"is_read"`
	ReadAt       *time.Time `db:"read_at" json:"read_at"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
