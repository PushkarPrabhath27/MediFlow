package request

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/equipment"
	"github.com/mediflow/backend/internal/shared/models"
)

type Service interface {
	CreateRequest(ctx context.Context, req *models.SharingRequest) error
	GetRequest(ctx context.Context, tenantID, reqID uuid.UUID) (*models.SharingRequest, error)
	ListRequests(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) ([]models.SharingRequest, error)
	ApproveRequest(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID) error
	DeclineRequest(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, reason string) error
	ConfirmHandoff(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, isSource bool) error
	ConfirmReturn(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, isSource bool) error
	AutoMatch(ctx context.Context, tenantID uuid.UUID, reqID uuid.UUID) error
}

type requestService struct {
	repo         Repository
	equipService equipment.Service
}

func NewService(repo Repository, equipService equipment.Service) Service {
	return &requestService{
		repo:         repo,
		equipService: equipService,
	}
}

func (s *requestService) CreateRequest(ctx context.Context, req *models.SharingRequest) error {
	req.Status = models.RequestStatusPending
	err := s.repo.Create(ctx, req)
	if err != nil {
		return err
	}

	// Trigger async auto-match
	go s.AutoMatch(context.Background(), req.TenantID, req.ID)
	
	return nil
}

func (s *requestService) GetRequest(ctx context.Context, tenantID, reqID uuid.UUID) (*models.SharingRequest, error) {
	return s.repo.FindByID(ctx, tenantID, reqID)
}

func (s *requestService) ListRequests(ctx context.Context, tenantID uuid.UUID, filters map[string]interface{}) ([]models.SharingRequest, error) {
	return s.repo.FindAll(ctx, tenantID, filters)
}

func (s *requestService) ApproveRequest(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID) error {
	req, err := s.repo.FindByID(ctx, tenantID, reqID)
	if err != nil {
		return err
	}

	if req.Status != models.RequestStatusMatched {
		return fmt.Errorf("request must be in matched state to approve")
	}

	now := time.Now()
	req.Status = models.RequestStatusApproved
	req.ApprovedAt = &now
	req.ApprovedByUserID = &userID

	return s.repo.Update(ctx, req)
}

func (s *requestService) DeclineRequest(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, reason string) error {
	return s.repo.UpdateStatus(ctx, reqID, models.RequestStatusDeclined, userID, reason)
}

func (s *requestService) ConfirmHandoff(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, isSource bool) error {
	err := s.repo.ConfirmHandoff(ctx, reqID, isSource)
	if err != nil {
		return err
	}

	// Re-fetch to check if both confirmed
	req, err := s.repo.FindByID(ctx, tenantID, reqID)
	if err != nil {
		return err
	}

	if req.HandoffConfirmedBySource && req.HandoffConfirmedByRequester {
		now := time.Now()
		req.Status = models.RequestStatusInTransit // Or Active if already moved
		req.HandedOffAt = &now
		
		// If equipment ID is set, update equipment status to InTransit
		if req.EquipmentID != nil {
			err = s.equipService.UpdateStatus(ctx, tenantID, *req.EquipmentID, models.StatusInTransit, nil, userID, "Shared handoff confirmed")
			if err != nil {
				return err
			}
		}

		return s.repo.Update(ctx, req)
	}

	return nil
}

func (s *requestService) ConfirmReturn(ctx context.Context, tenantID, reqID uuid.UUID, userID uuid.UUID, isSource bool) error {
	err := s.repo.ConfirmReturn(ctx, reqID, isSource)
	if err != nil {
		return err
	}

	req, err := s.repo.FindByID(ctx, tenantID, reqID)
	if err != nil {
		return err
	}

	if req.ReturnConfirmedBySource && req.ReturnConfirmedByRequester {
		now := time.Now()
		req.Status = models.RequestStatusCompleted
		req.ReturnedAt = &now

		if req.EquipmentID != nil {
			// Set back to available at source department
			err = s.equipService.UpdateStatus(ctx, tenantID, *req.EquipmentID, models.StatusAvailable, nil, userID, "Shared return completed")
			if err != nil {
				return err
			}
		}

		return s.repo.Update(ctx, req)
	}

	return nil
}

func (s *requestService) AutoMatch(ctx context.Context, tenantID uuid.UUID, reqID uuid.UUID) error {
	// Simple implementation: check all departments for availability
	// In production, this would be more efficient
	if _, err := s.repo.FindByID(ctx, tenantID, reqID); err != nil {
		return err
	}

	// This is a placeholder for the full matching logic
	// It would involve querying availability per department and scoring them
	fmt.Printf("Auto-matching request %s...\n", reqID)
	
	return nil
}
