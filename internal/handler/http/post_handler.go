package http

import (
	"net/http"
	"strconv"

	"moon/internal/domain/post"
	"moon/internal/usecase"
	"moon/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostHandler struct {
	postUseCase usecase.PostUseCase
	logger      *zap.Logger
}

// NewPostHandler creates a new post handler
func NewPostHandler(postUseCase usecase.PostUseCase) *PostHandler {
	return &PostHandler{
		postUseCase: postUseCase,
		logger:      logger.GetLogger(),
	}
}

// CreatePost handles creating a new post
// @Summary Create a new post
// @Description Create a new post (authenticated users)
// @Tags posts
// @Accept json
// @Produce json
// @Param request body post.CreatePostRequest true "Post creation data"
// @Success 201 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req post.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	postResponse, err := h.postUseCase.CreatePost(c.Request.Context(), req, userID.(uint))
	if err != nil {
		h.logger.Error("Failed to create post", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Post created successfully", zap.Uint("post_id", postResponse.ID), zap.Any("user_id", userID))
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data":    postResponse,
	})
}

// GetPostByID handles getting a post by ID
// @Summary Get post by ID
// @Description Get a specific post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param increment_view query bool false "Increment view count"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	incrementView := c.DefaultQuery("increment_view", "true") == "true"

	postResponse, err := h.postUseCase.GetPostByID(c.Request.Context(), uint(id), incrementView)
	if err != nil {
		h.logger.Error("Failed to get post", zap.Error(err), zap.Uint64("id", id))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved post", zap.Uint64("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post retrieved successfully",
		"data":    postResponse,
	})
}

// GetPostBySlug handles getting a post by slug
// @Summary Get post by slug
// @Description Get a specific post by slug
// @Tags posts
// @Accept json
// @Produce json
// @Param slug path string true "Post slug"
// @Param increment_view query bool false "Increment view count"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/slug/{slug} [get]
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		h.logger.Error("Empty post slug")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post slug is required",
		})
		return
	}

	incrementView := c.DefaultQuery("increment_view", "true") == "true"

	postResponse, err := h.postUseCase.GetPostBySlug(c.Request.Context(), slug, incrementView)
	if err != nil {
		h.logger.Error("Failed to get post by slug", zap.Error(err), zap.String("slug", slug))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved post by slug", zap.String("slug", slug))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post retrieved successfully",
		"data":    postResponse,
	})
}

// UpdatePost handles updating a post
// @Summary Update post
// @Description Update a post (author or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param request body post.UpdatePostRequest true "Post update data"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	var req post.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userRole, _ := c.Get("role")

	postResponse, err := h.postUseCase.UpdatePost(c.Request.Context(), uint(id), req, userID.(uint), userRole.(string))
	if err != nil {
		h.logger.Error("Failed to update post", zap.Error(err), zap.Uint64("id", id), zap.Any("user_id", userID))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Updated post", zap.Uint64("id", id), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"data":    postResponse,
	})
}

// DeletePost handles deleting a post
// @Summary Delete post
// @Description Delete a post (author or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userRole, _ := c.Get("role")

	err = h.postUseCase.DeletePost(c.Request.Context(), uint(id), userID.(uint), userRole.(string))
	if err != nil {
		h.logger.Error("Failed to delete post", zap.Error(err), zap.Uint64("id", id), zap.Any("user_id", userID))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Deleted post", zap.Uint64("id", id), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

// GetAllPosts handles getting all posts with filtering
// @Summary Get all posts
// @Description Get all posts with filtering and pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param status query string false "Post status" Enums(draft, published, archived)
// @Param category_id query int false "Category ID"
// @Param author_id query int false "Author ID"
// @Param is_public query bool false "Is public"
// @Param search query string false "Search in title and content"
// @Success 200 {object} post.PostsListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts [get]
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Build filter
	filter := post.PostFilter{}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryIDUint := uint(categoryID)
			filter.CategoryID = &categoryIDUint
		}
	}

	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if authorID, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			authorIDUint := uint(authorID)
			filter.AuthorID = &authorIDUint
		}
	}

	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filter.IsPublic = &isPublic
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	postsResponse, err := h.postUseCase.GetAllPosts(c.Request.Context(), filter, page, limit)
	if err != nil {
		h.logger.Error("Failed to get posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved posts list", zap.Int("count", len(postsResponse.Posts)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Posts retrieved successfully",
		"data":    postsResponse,
	})
}

// GetMyPosts handles getting current user's posts
// @Summary Get my posts
// @Description Get posts created by the authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} post.PostsListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/my [get]
func (h *PostHandler) GetMyPosts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	postsResponse, err := h.postUseCase.GetMyPosts(c.Request.Context(), userID.(uint), page, limit)
	if err != nil {
		h.logger.Error("Failed to get user posts", zap.Error(err), zap.Any("user_id", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved user posts", zap.Any("user_id", userID), zap.Int("count", len(postsResponse.Posts)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Posts retrieved successfully",
		"data":    postsResponse,
	})
}

// GetPublishedPosts handles getting published posts (public endpoint)
// @Summary Get published posts
// @Description Get all published and public posts
// @Tags posts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} post.PostsListResponse
// @Failure 500 {object} map[string]interface{}
// @Router /posts/published [get]
func (h *PostHandler) GetPublishedPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	postsResponse, err := h.postUseCase.GetPublishedPosts(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("Failed to get published posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Retrieved published posts", zap.Int("count", len(postsResponse.Posts)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Posts retrieved successfully",
		"data":    postsResponse,
	})
}

// PublishPost handles publishing a post
// @Summary Publish post
// @Description Publish a post (author or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/publish [patch]
func (h *PostHandler) PublishPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userRole, _ := c.Get("role")

	postResponse, err := h.postUseCase.PublishPost(c.Request.Context(), uint(id), userID.(uint), userRole.(string))
	if err != nil {
		h.logger.Error("Failed to publish post", zap.Error(err), zap.Uint64("id", id), zap.Any("user_id", userID))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Published post", zap.Uint64("id", id), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post published successfully",
		"data":    postResponse,
	})
}

// UnpublishPost handles unpublishing a post
// @Summary Unpublish post
// @Description Unpublish a post (author or admin only)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} post.PostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /posts/{id}/unpublish [patch]
func (h *PostHandler) UnpublishPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid post ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userRole, _ := c.Get("role")

	postResponse, err := h.postUseCase.UnpublishPost(c.Request.Context(), uint(id), userID.(uint), userRole.(string))
	if err != nil {
		h.logger.Error("Failed to unpublish post", zap.Error(err), zap.Uint64("id", id), zap.Any("user_id", userID))
		statusCode := http.StatusInternalServerError
		if err.Error() == "post not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "permission denied" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Unpublished post", zap.Uint64("id", id), zap.Any("user_id", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "Post unpublished successfully",
		"data":    postResponse,
	})
}
