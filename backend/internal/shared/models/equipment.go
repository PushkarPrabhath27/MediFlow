package models

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentStatus string

const (
	StatusAvailable     EquipmentStatus = "available"
	StatusInUse         EquipmentStatus = "in_use"
	StatusReserved      EquipmentStatus = "reserved"
	StatusInMaintenance EquipmentStatus = "in_maintenance"
	StatusInTransit     EquipmentStatus = "in_transit"
	StatusMissing       EquipmentStatus = "missing"
	StatusDecommissioned EquipmentStatus = "decommissioned"
)

type EquipmentItem struct {
	ID                uuid.UUID       `db:"id" json:"id"`
	TenantID          uuid.UUID       `db:"tenant_id" json:"tenant_id"`
	CategoryID        uuid.UUID       `db:"category_id" json:"category_id"`
	DepartmentID      *uuid.UUID      `db:"department_id" json:"department_id"`
	CurrentLocationID *uuid.UUID      `db:"current_location_id" json:"current_location_id"`
	Name              string          `db:"name" json:"name"`
	Model             string          `db:"model" json:"model"`
	Manufacturer      string          `db:"manufacturer" json:"manufacturer"`
	SerialNumber      string          `db:"serial_number" json:"serial_number"`
	AssetTag          string          `db:"asset_tag" json:"asset_tag"`
	QRCode            string          `db:"qr_code" json:"qr_code"`
	Status            EquipmentStatus `db:"status" json:"status"`
	PurchaseDate      *time.Time      `db:"purchase_date" json:"purchase_date"`
	PurchaseCost      float64         `db:"purchase_cost" json:"purchase_cost"`
	Notes             string          `db:"notes" json:"notes"`
	IsShared          bool            `db:"is_shared" json:"is_shared"`
	CreatedAt         time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
}

type EquipmentCategory struct {
	ID          uuid.UUID `db:"id" json:"id"`
	TenantID    uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Icon        string    `db:"icon" json:"icon"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Department struct {
	ID         uuid.UUID `db:"id" json:"id"`
	TenantID   uuid.UUID `db:"tenant_id" json:"tenant_id"`
	Name       string    `db:"name" json:"name"`
	Code       string    `db:"code" json:"code"`
	Floor      string    `db:"floor" json:"floor"`
	Building   string    `db:"building" json:"building"`
	HeadUserID *uuid.UUID `db:"head_user_id" json:"head_user_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type EquipmentStatusLog struct {
	ID               uuid.UUID       `db:"id" json:"id"`
	EquipmentID      uuid.UUID       `db:"equipment_id" json:"equipment_id"`
	TenantID         uuid.UUID       `db:"tenant_id" json:"tenant_id"`
	OldStatus        *EquipmentStatus `db:"old_status" json:"old_status"`
	NewStatus        EquipmentStatus `db:"new_status" json:"new_status"`
	OldDepartmentID  *uuid.UUID      `db:"old_department_id" json:"old_department_id"`
	NewDepartmentID  *uuid.UUID      `db:"new_department_id" json:"new_department_id"`
	OldLocationID    *uuid.UUID      `db:"old_location_id" json:"old_location_id"`
	NewLocationID    *uuid.UUID      `db:"new_location_id" json:"new_location_id"`
	ChangedByUserID  uuid.UUID       `db:"changed_by_user_id" json:"changed_by_user_id"`
	Reason           string          `db:"reason" json:"reason"`
	ChangedAt        time.Time       `db:"changed_at" json:"changed_at"`
}
