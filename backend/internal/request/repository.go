package request

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mediflow/backend/internal/shared/models"
)

type Repository interface {
	Create(ctx context.Context, req *models.SharingRequest) error
	FindByID(ctx context.Context, tenantID, reqID uuid.UUID) (*models.SharingRequest, error)
	FindAll(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) ([]models.SharingRequest, error)
	Update(ctx context.Context, req *models.SharingRequest) error
	UpdateStatus(ctx context.Context, reqID uuid.UUID, newStatus models.RequestStatus, userID uuid.UUID, notes string) error
	ConfirmHandoff(ctx context.Context, reqID uuid.UUID, isSource bool) error
	ConfirmReturn(ctx context.Context, reqID uuid.UUID, isSource bool) error
	CreateHistory(ctx context.Context, history *models.RequestHistory) error
}

type postgresRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, req *models.SharingRequest) error {
	query := `INSERT INTO sharing_requests (
		tenant_id, requesting_dept_id, requesting_user_id, category_id, 
		quantity_needed, urgency, reason, status, needed_by, expected_return_at
	) VALUES (
		:tenant_id, :requesting_dept_id, :requesting_user_id, :category_id, 
		:quantity_needed, :urgency, :reason, :status, :needed_by, :expected_return_at
	) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, req)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	}
	return fmt.Errorf("failed to create sharing request")
}

func (r *postgresRepo) FindByID(ctx context.Context, tenantID, reqID uuid.UUID) (*models.SharingRequest, error) {
	var req models.SharingRequest
	query := `SELECT * FROM sharing_requests WHERE tenant_id = $1 AND id = $2`
	err := r.db.GetContext(ctx, &req, query, tenantID, reqID)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *postgresRepo) FindAll(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) ([]models.SharingRequest, error) {
	var requests []models.SharingRequest
	query := `SELECT * FROM sharing_requests WHERE tenant_id = ?`
	args := []interface{}{tenantID}

	if deptID, ok := filters["dept_id"]; ok {
		query += " AND (requesting_dept_id = ? OR source_dept_id = ?)"
		args = append(args, deptID, deptID)
	}
	if status, ok := filters["status"]; ok {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"
	query = r.db.Rebind(query)
	err := r.db.SelectContext(ctx, &requests, query, args...)
	return requests, err
}

func (r *postgresRepo) Update(ctx context.Context, req *models.SharingRequest) error {
	query := `UPDATE sharing_requests SET 
		source_dept_id = :source_dept_id, equipment_id = :equipment_id,
		status = :status, matched_at = :matched_at, approved_at = :approved_at,
		approved_by_user_id = :approved_by_user_id, declined_reason = :declined_reason,
		handoff_confirmed_by_source = :handoff_confirmed_by_source,
		handoff_confirmed_by_requester = :handoff_confirmed_by_requester,
		handed_off_at = :handed_off_at,
		return_confirmed_by_source = :return_confirmed_by_source,
		return_confirmed_by_requester = :return_confirmed_by_requester,
		returned_at = :returned_at,
		updated_at = NOW()
	WHERE id = :id AND tenant_id = :tenant_id`

	_, err := r.db.NamedExecContext(ctx, query, req)
	return err
}

func (r *postgresRepo) UpdateStatus(ctx context.Context, reqID uuid.UUID, newStatus models.RequestStatus, userID uuid.UUID, notes string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldStatus models.RequestStatus
	err = tx.GetContext(ctx, &oldStatus, "SELECT status FROM sharing_requests WHERE id = $1 FOR UPDATE", reqID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE sharing_requests SET status = $1, updated_at = NOW() WHERE id = $2", newStatus, reqID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO request_history (request_id, changed_by_user_id, old_status, new_status, notes) VALUES ($1, $2, $3, $4, $5)",
		reqID, userID, oldStatus, newStatus, notes)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *postgresRepo) ConfirmHandoff(ctx context.Context, reqID uuid.UUID, isSource bool) error {
	field := "handoff_confirmed_by_requester"
	if isSource {
		field = "handoff_confirmed_by_source"
	}

	query := fmt.Sprintf("UPDATE sharing_requests SET %s = true, updated_at = NOW() WHERE id = $1", field)
	_, err := r.db.ExecContext(ctx, query, reqID)
	return err
}

func (r *postgresRepo) ConfirmReturn(ctx context.Context, reqID uuid.UUID, isSource bool) error {
	field := "return_confirmed_by_requester"
	if isSource {
		field = "return_confirmed_by_source"
	}

	query := fmt.Sprintf("UPDATE sharing_requests SET %s = true, updated_at = NOW() WHERE id = $1", field)
	_, err := r.db.ExecContext(ctx, query, reqID)
	return err
}

func (r *postgresRepo) CreateHistory(ctx context.Context, history *models.RequestHistory) error {
	query := `INSERT INTO request_history (request_id, changed_by_user_id, old_status, new_status, notes) 
	VALUES (:request_id, :changed_by_user_id, :old_status, :new_status, :notes)`
	_, err := r.db.NamedExecContext(ctx, query, history)
	return err
}
