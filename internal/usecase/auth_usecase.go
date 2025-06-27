package usecase

import (
	"context"
	"errors"

	"moon/internal/config"
	"moon/internal/domain/user"
	"moon/pkg/hash"
	"moon/pkg/jwt"
)

type AuthUseCase interface {
	Register(ctx context.Context, req user.CreateUserRequest) (*user.UserResponse, error)
	Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error)
}

type authUseCase struct {
	userRepo user.Repository
	cfg      *config.Config
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(userRepo user.Repository, cfg *config.Config) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (uc *authUseCase) Register(ctx context.Context, req user.CreateUserRequest) (*user.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	newUser := &user.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Phone:    nil,
		Address:  nil,
		Lat:      nil,
		Lng:      nil,
		Role:     "user",
		IsActive: true,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Return user response
	response := &user.UserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Name:      newUser.Name,
		Phone:     getStringValue(newUser.Phone),
		Address:   getStringValue(newUser.Address),
		Lat:       getFloat64Value(newUser.Lat),
		Lng:       getFloat64Value(newUser.Lng),
		Role:      newUser.Role,
		IsActive:  newUser.IsActive,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return response, nil
}

func (uc *authUseCase) Login(ctx context.Context, req user.LoginRequest) (*user.LoginResponse, error) {
	// Get user by email
	u, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !u.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Verify password
	if !hash.CheckPasswordHash(req.Password, u.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(u.ID, u.Email, u.Role, uc.cfg.JWT.Secret, uc.cfg.JWT.ExpiresIn)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Prepare user response
	userResponse := user.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Phone:     getStringValue(u.Phone),
		Address:   getStringValue(u.Address),
		Lat:       getFloat64Value(u.Lat),
		Lng:       getFloat64Value(u.Lng),
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	// Return login response
	return &user.LoginResponse{
		Token: token,
		User:  userResponse,
	}, nil
}

// Helper functions to handle pointer values
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func getFloat64Value(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}
