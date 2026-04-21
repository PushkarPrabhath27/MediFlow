package equipment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/shared/redis"
	"github.com/rs/zerolog/log"
)

type AvailabilityState struct {
	DepartmentID   uuid.UUID `json:"department_id"`
	CategoryID     uuid.UUID `json:"category_id"`
	AvailableCount int       `json:"available_count"`
	InUseCount     int       `json:"in_use_count"`
	TotalCount     int       `json:"total_count"`
}

type AvailabilityManager struct {
	repo        Repository
	redisClient *redis.Client
}

func NewAvailabilityManager(repo Repository, redisClient *redis.Client) *AvailabilityManager {
	return &AvailabilityManager{
		repo:        repo,
		redisClient: redisClient,
	}
}

// RefreshAll refresh the entire board for all tenants (should be run on startup/periodically)
func (m *AvailabilityManager) RefreshAll(ctx context.Context, tenantID uuid.UUID) error {
	summary, err := m.repo.GetAvailabilitySummary(ctx, tenantID)
	if err != nil {
		return err
	}

	boardKey := redis.GetAvailabilityBoardKey(tenantID.String())
	
	// Create a map for HSet
	values := make(map[string]interface{})
	for _, row := range summary {
		deptID := row["department_id"].(uuid.UUID).String()
		catID := row["category_id"].(uuid.UUID).String()
		field := fmt.Sprintf("%s:%s", deptID, catID)
		
		data, _ := json.Marshal(row)
		values[field] = data
	}

	if len(values) > 0 {
		err = m.redisClient.HSet(ctx, boardKey, values).Err()
		if err != nil {
			return err
		}
		m.redisClient.Expire(ctx, boardKey, 90*time.Second)
	}

	return nil
}

// UpdateCell updates a specific cell in Redis and publishes the event
func (m *AvailabilityManager) UpdateCell(ctx context.Context, tenantID, departmentID, categoryID uuid.UUID) error {
	// 1. Query latest from DB (could be optimized)
	summary, err := m.repo.GetAvailabilitySummary(ctx, tenantID)
	if err != nil {
		return err
	}

	var targetRow map[string]interface{}
	for _, row := range summary {
		if row["department_id"].(uuid.UUID) == departmentID && row["category_id"].(uuid.UUID) == categoryID {
			targetRow = row
			break
		}
	}

	if targetRow == nil {
		return nil // Nothing to update
	}

	// 2. Update Redis Hash
	boardKey := redis.GetAvailabilityBoardKey(tenantID.String())
	field := fmt.Sprintf("%s:%s", departmentID.String(), categoryID.String())
	data, _ := json.Marshal(targetRow)
	
	err = m.redisClient.HSet(ctx, boardKey, field, data).Err()
	if err != nil {
		log.Error().Err(err).Msg("Failed to update Redis availability hash")
	}

	// 3. Publish to Redis for WebSocket bridge
	publishChannel := fmt.Sprintf("ws-events:%s", tenantID.String())
	event := map[string]interface{}{
		"type": "availability_update",
		"data": targetRow,
	}
	eventData, _ := json.Marshal(event)
	
	return m.redisClient.Publish(ctx, publishChannel, eventData).Err()
}
