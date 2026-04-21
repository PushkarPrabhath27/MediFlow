package analytics

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/shared/middleware"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/analytics", func(r chi.Router) {
		r.Get("/dashboard", h.GetDashboard)
	})
}

// @Summary Get dashboard analytics
// @Tags analytics
// @Produce json
// @Success 200 {array} models.UtilizationReport
// @Router /api/v1/analytics/dashboard [get]
func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	// If user is department head or lower, restrict to their department
	// In a real app, this would use a more robust permission check
	var deptID *uuid.UUID
	// Mock logic: if query param dept_id is provided, use it
	if dID := r.URL.Query().Get("dept_id"); dID != "" {
		parsedID, err := uuid.Parse(dID)
		if err == nil {
			deptID = &parsedID
		}
	}

	reports, err := h.service.GetDashboardMetrics(r.Context(), claims.TenantID, deptID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard metrics")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reports)
}
