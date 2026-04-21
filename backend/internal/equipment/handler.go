package equipment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mediflow/backend/internal/shared/middleware"
	"github.com/mediflow/backend/internal/shared/models"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/equipment", func(r chi.Router) {
		r.Get("/", h.ListItems)
		r.Post("/", h.CreateItem)
		r.Get("/availability-summary", h.GetAvailabilitySummary)
		r.Post("/qr-update", h.UpdateStatusByQR)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetItem)
			r.Put("/status", h.UpdateStatus)
			r.Get("/history", h.GetStatusHistory)
		})
	})
}

// @Summary List equipment items
// @Tags equipment
// @Produce json
// @Param category_id query string false "Category ID"
// @Param department_id query string false "Department ID"
// @Param status query string false "Status"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.EquipmentItem
// @Router /api/v1/equipment [get]
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	
	filters := make(map[string]interface{})
	if catID := r.URL.Query().Get("category_id"); catID != "" {
		filters["category_id"] = catID
	}
	if deptID := r.URL.Query().Get("department_id"); deptID != "" {
		filters["department_id"] = deptID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	items, err := h.service.ListItems(r.Context(), claims.TenantID, filters, offset, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list items")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

// @Summary Create equipment item
// @Tags equipment
// @Accept json
// @Produce json
// @Param item body models.EquipmentItem true "Item details"
// @Success 201 {object} models.EquipmentItem
// @Router /api/v1/equipment [post]
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	var item models.EquipmentItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item.TenantID = claims.TenantID
	
	if err := h.service.CreateItem(r.Context(), &item); err != nil {
		log.Error().Err(err).Msg("Failed to create item")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// @Summary Get equipment item
// @Tags equipment
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.EquipmentItem
// @Router /api/v1/equipment/{id} [get]
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetItem(r.Context(), claims.TenantID, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get item")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}

// @Summary Update equipment status
// @Tags equipment
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param request body UpdateStatusRequest true "Status update details"
// @Success 200 {string} string "Status updated"
// @Router /api/v1/equipment/{id}/status [put]
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	type UpdateStatusRequest struct {
		NewStatus  models.EquipmentStatus `json:"new_status"`
		LocationID *uuid.UUID             `json:"location_id"`
		Reason     string                 `json:"reason"`
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateStatus(r.Context(), claims.TenantID, itemID, req.NewStatus, req.LocationID, claims.UserID, req.Reason)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update status")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "status updated"}`))
}

// @Summary Update equipment status by QR code
// @Tags equipment
// @Accept json
// @Produce json
// @Param request body UpdateStatusByQRRequest true "QR status update details"
// @Success 200 {string} string "Status updated"
// @Router /api/v1/equipment/qr-update [post]
func (h *Handler) UpdateStatusByQR(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	type UpdateStatusByQRRequest struct {
		QRCode     string                 `json:"qr_code"`
		NewStatus  models.EquipmentStatus `json:"new_status"`
		LocationID *uuid.UUID             `json:"location_id"`
		Reason     string                 `json:"reason"`
	}

	var req UpdateStatusByQRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateStatusByQR(r.Context(), req.QRCode, req.NewStatus, req.LocationID, claims.UserID, req.Reason)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update status by QR")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "status updated"}`))
}

// @Summary Get equipment status history
// @Tags equipment
// @Produce json
// @Param id path string true "Item ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.EquipmentStatusLog
// @Router /api/v1/equipment/{id}/history [get]
func (h *Handler) GetStatusHistory(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	logs, err := h.service.GetStatusHistory(r.Context(), claims.TenantID, itemID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get history")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(logs)
}

// @Summary Get availability summary
// @Tags equipment
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/equipment/availability-summary [get]
func (h *Handler) GetAvailabilitySummary(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	summary, err := h.service.GetAvailabilitySummary(r.Context(), claims.TenantID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get summary")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(summary)
}
