package alert

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mediflow/backend/internal/shared/models"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Service interface {
	NotifyUser(ctx context.Context, tenantID, userID uuid.UUID, nType, title, message string, data interface{}) error
	NotifyDepartment(ctx context.Context, tenantID, deptID uuid.UUID, nType, title, message string, data interface{}) error
	ListNotifications(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, tenantID, notificationID uuid.UUID) error
}

type alertService struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func NewService(db *sqlx.DB, redisClient *redis.Client) Service {
	return &alertService{
		db:          db,
		redisClient: redisClient,
	}
}

func (s *alertService) NotifyUser(ctx context.Context, tenantID, userID uuid.UUID, nType, title, message string, data interface{}) error {
	dataJSON, _ := json.Marshal(data)
	
	query := `INSERT INTO notifications (tenant_id, user_id, type, title, message, data_json) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, tenantID, userID, nType, title, message, dataJSON)
	if err != nil {
		return err
	}

	// Publish for real-time WebSocket delivery
	return s.publishRealtime(ctx, tenantID, "user_notification", map[string]interface{}{
		"user_id": userID,
		"type":    nType,
		"title":   title,
		"message": message,
		"data":    data,
	})
}

func (s *alertService) NotifyDepartment(ctx context.Context, tenantID, deptID uuid.UUID, nType, title, message string, data interface{}) error {
	dataJSON, _ := json.Marshal(data)
	
	query := `INSERT INTO notifications (tenant_id, department_id, type, title, message, data_json) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, query, tenantID, deptID, nType, title, message, dataJSON)
	if err != nil {
		return err
	}

	// Publish for real-time WebSocket delivery
	return s.publishRealtime(ctx, tenantID, "dept_notification", map[string]interface{}{
		"department_id": deptID,
		"type":          nType,
		"title":         title,
		"message":       message,
		"data":          data,
	})
}

func (s *alertService) ListNotifications(ctx context.Context, tenantID, userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var results []models.Notification
	query := `SELECT * FROM notifications WHERE tenant_id = $1 AND (user_id = $2 OR department_id IN (SELECT department_id FROM users WHERE id = $2)) 
	ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	err := s.db.SelectContext(ctx, &results, query, tenantID, userID, limit, offset)
	return results, err
}

func (s *alertService) MarkAsRead(ctx context.Context, tenantID, notificationID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, "UPDATE notifications SET is_read = true, read_at = NOW() WHERE id = $1 AND tenant_id = $2", notificationID, tenantID)
	return err
}

func (s *alertService) publishRealtime(ctx context.Context, tenantID uuid.UUID, eType string, payload interface{}) error {
	channel := fmt.Sprintf("ws-events:%s", tenantID.String())
	event := map[string]interface{}{
		"type":    eType,
		"payload": payload,
	}
	data, _ := json.Marshal(event)
	
	err := s.redisClient.Publish(ctx, channel, data).Err()
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish real-time notification")
	}
	return err
}
