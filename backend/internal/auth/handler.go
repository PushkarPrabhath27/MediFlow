package auth

import (
	"encoding/json"
	"net/http"

	"github.com/mediflow/backend/internal/shared/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type Handler struct {
	jwtManager *JWTManager
}

func NewHandler(jwtManager *JWTManager) *Handler {
	return &Handler{jwtManager: jwtManager}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In Phase 1, we use a mock authentication for testing
	// Real implementation would check DB password hash
	mockUser := &models.User{
		Email: req.Email,
		Role:  models.RoleStaff, // Default
	}
	
	// Simple mock role based on email for testing
	if req.Email == "admin@hospital.com" {
		mockUser.Role = models.RoleHospitalAdmin
	}

	token, err := h.jwtManager.Generate(mockUser)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  mockUser,
	})
}
