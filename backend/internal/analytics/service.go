package analytics

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/shared/models"
)

type Service interface {
	GetDashboardMetrics(ctx context.Context, tenantID uuid.UUID, deptID *uuid.UUID) ([]models.UtilizationReport, error)
}

type analyticsService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &analyticsService{
		repo: repo,
	}
}

func (s *analyticsService) GetDashboardMetrics(ctx context.Context, tenantID uuid.UUID, deptID *uuid.UUID) ([]models.UtilizationReport, error) {
	// For dashboard, get last 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	return s.repo.GetUtilizationReport(ctx, tenantID, startDate, endDate, deptID)
}
