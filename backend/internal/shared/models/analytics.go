package models

import (
	"time"

	"github.com/google/uuid"
)

type UtilizationMetric struct {
	ID                       uuid.UUID `db:"id" json:"id"`
	TenantID                 uuid.UUID `db:"tenant_id" json:"tenant_id"`
	DepartmentID             uuid.UUID `db:"department_id" json:"department_id"`
	CategoryID               uuid.UUID `db:"category_id" json:"category_id"`
	Date                     time.Time `db:"date" json:"date"`
	TotalHoursAvailable      float64   `db:"total_hours_available" json:"total_hours_available"`
	TotalHoursInUse          float64   `db:"total_hours_in_use" json:"total_hours_in_use"`
	TotalHoursInMaintenance  float64   `db:"total_hours_in_maintenance" json:"total_hours_in_maintenance"`
	SharingRequestsSent      int       `db:"sharing_requests_sent" json:"sharing_requests_sent"`
	SharingRequestsReceived  int       `db:"sharing_requests_received" json:"sharing_requests_received"`
	SharingRequestsFulfilled int       `db:"sharing_requests_fulfilled" json:"sharing_requests_fulfilled"`
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
}

type UtilizationReport struct {
	CategoryID           uuid.UUID `json:"category_id"`
	CategoryName         string    `json:"category_name"`
	TotalItems           int       `json:"total_items"`
	UtilizationRate      float64   `json:"utilization_rate"` // Percentage
	AvgTimeInMaintenance float64   `json:"avg_time_in_maintenance"`
	SharingSuccessRate   float64   `json:"sharing_success_rate"` // Percentage
}
