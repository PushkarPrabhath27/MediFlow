package request

import (
	"encoding/json"
	"net/http"

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
	r.Route("/requests", func(r chi.Router) {
		r.Get("/", h.ListRequests)
		r.Post("/", h.CreateRequest)
		
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetRequest)
			r.Post("/approve", h.ApproveRequest)
			r.Post("/decline", h.DeclineRequest)
			r.Post("/confirm-handoff", h.ConfirmHandoff)
			r.Post("/confirm-return", h.ConfirmReturn)
		})
	})
}

func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	var req models.SharingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	req.TenantID = claims.TenantID
	req.RequestingUserID = claims.UserID
	
	if err := h.service.CreateRequest(r.Context(), &req); err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *Handler) ListRequests(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())

	filters := make(map[string]interface{})
	if deptID := r.URL.Query().Get("dept_id"); deptID != "" {
		filters["dept_id"] = deptID
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	requests, err := h.service.ListRequests(r.Context(), claims.TenantID, filters)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list requests")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	req, err := h.service.GetRequest(r.Context(), claims.TenantID, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get request")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(req)
}

func (h *Handler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)

	err := h.service.ApproveRequest(r.Context(), claims.TenantID, id, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "approved"}`))
}

func (h *Handler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)

	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := h.service.DeclineRequest(r.Context(), claims.TenantID, id, claims.UserID, body.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "declined"}`))
}

func (h *Handler) ConfirmHandoff(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)

	var body struct {
		IsSource bool `json:"is_source"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := h.service.ConfirmHandoff(r.Context(), claims.TenantID, id, claims.UserID, body.IsSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "handoff confirmed"}`))
}

func (h *Handler) ConfirmReturn(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)

	var body struct {
		IsSource bool `json:"is_source"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := h.service.ConfirmReturn(r.Context(), claims.TenantID, id, claims.UserID, body.IsSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "return confirmed"}`))
}
