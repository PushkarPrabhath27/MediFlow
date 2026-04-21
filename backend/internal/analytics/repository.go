package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mediflow/backend/internal/shared/models"
)

type Repository interface {
	GetUtilizationReport(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, deptID *uuid.UUID) ([]models.UtilizationReport, error)
	TrackSharingRequest(ctx context.Context, tenantID, deptID, catID uuid.UUID, date time.Time, isSent, isFulfilled bool) error
}

type postgresRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetUtilizationReport(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time, deptID *uuid.UUID) ([]models.UtilizationReport, error) {
	// Simple mock implementation for Phase 4 to aggregate metrics
	// In a real application, this would join utilization_metrics, equipment_categories, and equipment_items
	var reports []models.UtilizationReport

	query := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			COUNT(e.id) as total_items,
			COALESCE(SUM(u.total_hours_in_use) / NULLIF(SUM(u.total_hours_available + u.total_hours_in_use + u.total_hours_in_maintenance), 0) * 100, 0) as utilization_rate,
			COALESCE(SUM(u.total_hours_in_maintenance) / COUNT(e.id), 0) as avg_time_in_maintenance,
			COALESCE(CAST(SUM(u.sharing_requests_fulfilled) AS FLOAT) / NULLIF(SUM(u.sharing_requests_received), 0) * 100, 0) as sharing_success_rate
		FROM equipment_categories c
		LEFT JOIN equipment_items e ON e.category_id = c.id
		LEFT JOIN utilization_metrics u ON u.category_id = c.id AND u.date >= $2 AND u.date <= $3
		WHERE c.tenant_id = $1
	`
	args := []interface{}{tenantID, startDate, endDate}

	if deptID != nil {
		query += ` AND (u.department_id = $4 OR u.department_id IS NULL) AND (e.department_id = $4)`
		args = append(args, *deptID)
	}

	query += ` GROUP BY c.id, c.name`

	err := r.db.SelectContext(ctx, &reports, query, args...)
	return reports, err
}

func (r *postgresRepo) TrackSharingRequest(ctx context.Context, tenantID, deptID, catID uuid.UUID, date time.Time, isSent, isFulfilled bool) error {
	sentInc := 0
	receivedInc := 0
	if isSent {
		sentInc = 1
	} else {
		receivedInc = 1
	}

	fulfilledInc := 0
	if isFulfilled {
		fulfilledInc = 1
	}

	query := `
		INSERT INTO utilization_metrics (
			tenant_id, department_id, category_id, date, 
			sharing_requests_sent, sharing_requests_received, sharing_requests_fulfilled
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, department_id, category_id, date) DO UPDATE SET
			sharing_requests_sent = utilization_metrics.sharing_requests_sent + $5,
			sharing_requests_received = utilization_metrics.sharing_requests_received + $6,
			sharing_requests_fulfilled = utilization_metrics.sharing_requests_fulfilled + $7,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, tenantID, deptID, catID, date, sentInc, receivedInc, fulfilledInc)
	return err
}
