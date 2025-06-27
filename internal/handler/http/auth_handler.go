package http

import (
	"net/http"

	"moon/internal/domain/user"
	"moon/internal/usecase"
	"moon/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	logger      *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		logger:      logger.GetLogger(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "User registration data"
// @Success 201 {object} user.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userResponse, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err), zap.String("email", req.Email))
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("User registered successfully", zap.String("email", req.Email), zap.Uint("user_id", userResponse.ID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    userResponse,
	})
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "User login credentials"
// @Success 200 {object} user.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	loginResponse, err := h.authUseCase.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err), zap.String("email", req.Email))
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid email or password" || err.Error() == "user account is deactivated" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("User logged in successfully", zap.String("email", req.Email), zap.Uint("user_id", loginResponse.User.ID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    loginResponse,
	})
}

// RefreshToken handles token refresh (optional - can be implemented later)
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// TODO: Implement refresh token logic
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Refresh token not implemented yet",
	})
}

// Logout handles user logout (optional - for token blacklisting)
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement logout logic (token blacklisting)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
