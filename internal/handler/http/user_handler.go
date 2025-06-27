package http

import (
	"net/http"
	"strconv"

	"moon/internal/domain/user"
	"moon/internal/usecase"
	"moon/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger.GetLogger(),
	}
}

// GetAllUsers handles getting all users (admin only)
// @Summary Get all users
// @Description Get all users with pagination (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} user.UsersListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	usersResponse, err := h.userUseCase.GetAllUsers(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved users list", zap.Int("count", len(usersResponse.Users)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    usersResponse,
	})
}

// GetUserByID handles getting a user by ID (admin only)
// @Summary Get user by ID
// @Description Get a specific user by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	userResponse, err := h.userUseCase.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err), zap.Uint64("id", id))
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved user", zap.Uint64("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    userResponse,
	})
}

// UpdateUser handles updating a user (admin only)
// @Summary Update user
// @Description Update a user's information (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body user.AdminUpdateUserRequest true "User update data"
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req user.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userResponse, err := h.userUseCase.UpdateUser(c.Request.Context(), uint(id), req)
	if err != nil {
		h.logger.Error("Failed to update user", zap.Error(err), zap.Uint64("id", id))
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Updated user", zap.Uint64("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    userResponse,
	})
}

// DeleteUser handles deleting a user (admin only)
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Prevent admin from deleting themselves
	currentUserID, _ := c.Get("user_id")
	if currentUserID == uint(id) {
		h.logger.Warn("Admin tried to delete themselves", zap.Uint64("id", id))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete your own account",
		})
		return
	}

	err = h.userUseCase.DeleteUser(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err), zap.Uint64("id", id))
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Deleted user", zap.Uint64("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetUsersByRole handles getting users by role (admin only)
// @Summary Get users by role
// @Description Get users filtered by role with pagination (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param role path string true "User role" Enums(user, admin)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} user.UsersListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users/role/{role} [get]
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
	role := c.Param("role")
	if role != "user" && role != "admin" {
		h.logger.Error("Invalid role", zap.String("role", role))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role. Must be 'user' or 'admin'",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	usersResponse, err := h.userUseCase.GetUsersByRole(c.Request.Context(), role, page, limit)
	if err != nil {
		h.logger.Error("Failed to get users by role", zap.Error(err), zap.String("role", role))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved users by role", zap.String("role", role), zap.Int("count", len(usersResponse.Users)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    usersResponse,
	})
}

// GetProfile handles getting current user profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} user.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userResponse, err := h.userUseCase.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err), zap.Any("user_id", userID))
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved user profile", zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    userResponse,
	})
}
