package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rmarko/electronics-marketplace/backend/internal/model"
	"github.com/rmarko/electronics-marketplace/backend/internal/service"
)

type ctxKeyUser struct{}

type AuthUser struct {
	UserID int
	Email  string
	Role   model.Role
}

func GetUser(ctx context.Context) *AuthUser {
	if u, ok := ctx.Value(ctxKeyUser{}).(*AuthUser); ok {
		return u
	}
	return nil
}

// RequireAuth rejects requests without a valid JWT.
func RequireAuth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				model.WriteError(w, model.ErrUnauthorized("missing authorization header"))
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				model.WriteError(w, model.ErrUnauthorized("invalid authorization format"))
				return
			}

			claims, err := authSvc.ValidateToken(parts[1])
			if err != nil {
				model.WriteError(w, err)
				return
			}

			user := &AuthUser{
				UserID: claims.UserID,
				Email:  claims.Email,
				Role:   claims.Role,
			}
			ctx := context.WithValue(r.Context(), ctxKeyUser{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireSeller rejects requests from non-seller users. Must be used after RequireAuth.
func RequireSeller(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil || user.Role != model.RoleSeller {
			model.WriteError(w, model.ErrForbidden("seller access required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireBuyer rejects requests from non-buyer users. Must be used after RequireAuth.
func RequireBuyer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		if user == nil || user.Role != model.RoleBuyer {
			model.WriteError(w, model.ErrForbidden("buyer access required"))
			return
		}
		next.ServeHTTP(w, r)
	})
}
