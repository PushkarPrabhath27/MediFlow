package equipment

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mediflow/backend/internal/shared/models"
)

type Repository interface {
	Create(ctx context.Context, item *models.EquipmentItem) error
	FindByID(ctx context.Context, tenantID, itemID uuid.UUID) (*models.EquipmentItem, error)
	FindAll(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}, offset, limit int) ([]models.EquipmentItem, error)
	Update(ctx context.Context, item *models.EquipmentItem) error
	UpdateStatus(ctx context.Context, itemID uuid.UUID, newStatus models.EquipmentStatus, locationID *uuid.UUID, changedByUserID uuid.UUID, reason string) error
	FindByQRCode(ctx context.Context, qrCode string) (*models.EquipmentItem, error)
	GetAvailabilitySummary(ctx context.Context, tenantID uuid.UUID) ([]map[string]interface{}, error)
	CreateStatusLog(ctx context.Context, log *models.EquipmentStatusLog) error
	FindStatusHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]models.EquipmentStatusLog, error)
}

type postgresRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, item *models.EquipmentItem) error {
	query := `INSERT INTO equipment_items (
		tenant_id, category_id, department_id, current_location_id, name, model, 
		manufacturer, serial_number, asset_tag, qr_code, status, purchase_date, 
		purchase_cost, notes, is_shared
	) VALUES (
		:tenant_id, :category_id, :department_id, :current_location_id, :name, :model, 
		:manufacturer, :serial_number, :asset_tag, :qr_code, :status, :purchase_date, 
		:purchase_cost, :notes, :is_shared
	) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	}
	return fmt.Errorf("failed to create equipment item")
}

func (r *postgresRepo) FindByID(ctx context.Context, tenantID, itemID uuid.UUID) (*models.EquipmentItem, error) {
	var item models.EquipmentItem
	query := `SELECT * FROM equipment_items WHERE tenant_id = $1 AND id = $2`
	err := r.db.GetContext(ctx, &item, query, tenantID, itemID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepo) FindAll(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}, offset, limit int) ([]models.EquipmentItem, error) {
	var items []models.EquipmentItem
	query := `SELECT * FROM equipment_items WHERE tenant_id = ?`
	args := []interface{}{tenantID}

	// Dynamic filtering
	if catID, ok := filters["category_id"]; ok {
		query += " AND category_id = ?"
		args = append(args, catID)
	}
	if deptID, ok := filters["department_id"]; ok {
		query += " AND department_id = ?"
		args = append(args, deptID)
	}
	if status, ok := filters["status"]; ok {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	query = r.db.Rebind(query)
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

func (r *postgresRepo) Update(ctx context.Context, item *models.EquipmentItem) error {
	query := `UPDATE equipment_items SET 
		category_id = :category_id, department_id = :department_id, 
		current_location_id = :current_location_id, name = :name, model = :model, 
		manufacturer = :manufacturer, serial_number = :serial_number, 
		asset_tag = :asset_tag, status = :status, purchase_date = :purchase_date, 
		purchase_cost = :purchase_cost, notes = :notes, is_shared = :is_shared,
		updated_at = NOW()
	WHERE id = :id AND tenant_id = :tenant_id`

	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *postgresRepo) UpdateStatus(ctx context.Context, itemID uuid.UUID, newStatus models.EquipmentStatus, locationID *uuid.UUID, changedByUserID uuid.UUID, reason string) error {
	// Start transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Get current status and details
	var current models.EquipmentItem
	err = tx.GetContext(ctx, &current, "SELECT status, tenant_id, department_id, current_location_id FROM equipment_items WHERE id = $1 FOR UPDATE", itemID)
	if err != nil {
		return err
	}

	// 2. Update status
	query := `UPDATE equipment_items SET status = $1, current_location_id = $2, updated_at = NOW() WHERE id = $3`
	locID := current.CurrentLocationID
	if locationID != nil {
		locID = locationID
	}
	_, err = tx.ExecContext(ctx, query, newStatus, locID, itemID)
	if err != nil {
		return err
	}

	// 3. Create log
	logQuery := `INSERT INTO equipment_status_logs (
		equipment_id, tenant_id, old_status, new_status, 
		old_department_id, new_department_id, 
		old_location_id, new_location_id, 
		changed_by_user_id, reason
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	
	_, err = tx.ExecContext(ctx, logQuery, 
		itemID, current.TenantID, current.Status, newStatus,
		current.DepartmentID, current.DepartmentID, // In Phase 1, department changes are handled separately
		current.CurrentLocationID, locID,
		changedByUserID, reason,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *postgresRepo) FindByQRCode(ctx context.Context, qrCode string) (*models.EquipmentItem, error) {
	var item models.EquipmentItem
	err := r.db.GetContext(ctx, &item, "SELECT * FROM equipment_items WHERE qr_code = $1", qrCode)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *postgresRepo) GetAvailabilitySummary(ctx context.Context, tenantID uuid.UUID) ([]map[string]interface{}, error) {
	// Aggregate counts per department per category
	query := `
		SELECT 
			d.id as department_id, 
			d.name as department_name,
			c.id as category_id, 
			c.name as category_name,
			COUNT(e.id) FILTER (WHERE e.status = 'available') as available_count,
			COUNT(e.id) FILTER (WHERE e.status = 'in_use') as in_use_count,
			COUNT(e.id) as total_count
		FROM departments d
		CROSS JOIN equipment_categories c
		LEFT JOIN equipment_items e ON e.department_id = d.id AND e.category_id = c.id
		WHERE d.tenant_id = $1 AND c.tenant_id = $1
		GROUP BY d.id, d.name, c.id, c.name
	`
	rows, err := r.db.QueryxContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		m := make(map[string]interface{})
		err := rows.MapScan(m)
		if err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (r *postgresRepo) CreateStatusLog(ctx context.Context, log *models.EquipmentStatusLog) error {
	query := `INSERT INTO equipment_status_logs (
		equipment_id, tenant_id, old_status, new_status, 
		old_department_id, new_department_id, 
		old_location_id, new_location_id, 
		changed_by_user_id, reason
	) VALUES (
		:equipment_id, :tenant_id, :old_status, :new_status, 
		:old_department_id, :new_department_id, 
		:old_location_id, :new_location_id, 
		:changed_by_user_id, :reason
	)`
	_, err := r.db.NamedExecContext(ctx, query, log)
	return err
}

func (r *postgresRepo) FindStatusHistory(ctx context.Context, itemID uuid.UUID, limit, offset int) ([]models.EquipmentStatusLog, error) {
	var logs []models.EquipmentStatusLog
	query := `SELECT * FROM equipment_status_logs WHERE equipment_id = $1 ORDER BY changed_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &logs, query, itemID, limit, offset)
	return logs, err
}
