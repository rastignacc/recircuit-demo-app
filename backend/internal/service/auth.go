package service

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/model"
	"github.com/rastignacc/recircuit-demo-app/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users     repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(users repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: []byte(jwtSecret)}
}

type Claims struct {
	UserID int        `json:"user_id"`
	Email  string     `json:"email"`
	Role   model.Role `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, model.ErrBadRequest("email, password, and name are required")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, model.ErrBadRequest("invalid email format")
	}
	if req.Role != model.RoleBuyer && req.Role != model.RoleSeller {
		return nil, model.ErrBadRequest("role must be 'buyer' or 'seller'")
	}
	if len(req.Password) < 6 {
		return nil, model.ErrBadRequest("password must be at least 6 characters")
	}

	existing, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrInternal("failed to check existing user")
	}
	if existing != nil {
		return nil, model.ErrConflict("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, model.ErrInternal("failed to hash password")
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         req.Role,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, model.ErrInternal("failed to create user")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, model.ErrInternal("failed to generate token")
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, model.ErrBadRequest("email and password are required")
	}

	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUnauthorized("invalid credentials")
		}
		return nil, model.ErrInternal("failed to find user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, model.ErrUnauthorized("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, model.ErrInternal("failed to generate token")
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrUnauthorized("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, model.ErrUnauthorized("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, model.ErrUnauthorized("invalid token claims")
	}
	return claims, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
