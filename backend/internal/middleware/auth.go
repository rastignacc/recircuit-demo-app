package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rastignacc/recircuit-demo-app/backend/internal/model"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/service"
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
// It checks the "token" HttpOnly cookie first, then falls back to the
// Authorization: Bearer header for API consumers.
func RequireAuth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""
			if c, err := r.Cookie("token"); err == nil && c.Value != "" {
				tokenStr = c.Value
			}
			if tokenStr == "" {
				header := r.Header.Get("Authorization")
				if header == "" {
					model.WriteError(w, model.ErrUnauthorized("missing authorization"))
					return
				}
				parts := strings.SplitN(header, " ", 2)
				if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
					model.WriteError(w, model.ErrUnauthorized("invalid authorization format"))
					return
				}
				tokenStr = parts[1]
			}

			claims, err := authSvc.ValidateToken(tokenStr)
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
