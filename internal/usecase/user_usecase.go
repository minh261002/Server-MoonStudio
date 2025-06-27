package usecase

import (
	"context"
	"errors"
	"math"

	"moon/internal/domain/user"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, page, limit int) (*user.UsersListResponse, error)
	GetUserByID(ctx context.Context, id uint) (*user.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req user.AdminUpdateUserRequest) (*user.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
	GetUsersByRole(ctx context.Context, role string, page, limit int) (*user.UsersListResponse, error)
}

type userUseCase struct {
	userRepo user.Repository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo user.Repository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (uc *userUseCase) GetAllUsers(ctx context.Context, page, limit int) (*user.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := uc.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	total, err := uc.userRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, errors.New("failed to count users")
	}

	userResponses := make([]user.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = user.UserResponse{
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
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &user.UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id uint) (*user.UserResponse, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user.UserResponse{
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
	}, nil
}

func (uc *userUseCase) UpdateUser(ctx context.Context, id uint, req user.AdminUpdateUserRequest) (*user.UserResponse, error) {
	u, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Phone != nil {
		u.Phone = req.Phone
	}
	if req.Address != nil {
		u.Address = req.Address
	}
	if req.Lat != nil {
		u.Lat = req.Lat
	}
	if req.Lng != nil {
		u.Lng = req.Lng
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
	if req.Role != nil {
		u.Role = *req.Role
	}

	if err := uc.userRepo.Update(ctx, u); err != nil {
		return nil, errors.New("failed to update user")
	}

	return &user.UserResponse{
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
	}, nil
}

func (uc *userUseCase) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (uc *userUseCase) GetUsersByRole(ctx context.Context, role string, page, limit int) (*user.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := uc.userRepo.GetByRole(ctx, role, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch users by role")
	}

	// Count users by role (you might want to add this method to repository)
	total, err := uc.userRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, errors.New("failed to count users")
	}

	userResponses := make([]user.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = user.UserResponse{
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
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &user.UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
