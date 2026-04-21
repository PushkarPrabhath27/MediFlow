package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mediflow/backend/internal/auth"
	"github.com/mediflow/backend/internal/shared/models"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	UserContextKey contextKey = "user_claims"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Also check query param for WebSockets
				tokenParam := r.URL.Query().Get("token")
				if tokenParam != "" {
					authHeader = "Bearer " + tokenParam
				}
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.Verify(tokenStr)
			if err != nil {
				log.Error().Err(err).Msg("Invalid token")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			isAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					isAllowed = true
					break
				}
			}

			// Super Admin always allowed
			if claims.Role == models.RoleSuperAdmin {
				isAllowed = true
			}

			if !isAllowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims retrieves claims from context
func GetClaims(ctx context.Context) *auth.Claims {
	claims, _ := ctx.Value(UserContextKey).(*auth.Claims)
	return claims
}
