package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"moon/internal/domain/post"
	"moon/internal/domain/user"
)

type PostUseCase interface {
	CreatePost(ctx context.Context, req post.CreatePostRequest, authorID uint) (*post.PostResponse, error)
	GetPostByID(ctx context.Context, id uint, incrementView bool) (*post.PostResponse, error)
	GetPostBySlug(ctx context.Context, slug string, incrementView bool) (*post.PostResponse, error)
	UpdatePost(ctx context.Context, id uint, req post.UpdatePostRequest, userID uint, userRole string) (*post.PostResponse, error)
	DeletePost(ctx context.Context, id uint, userID uint, userRole string) error
	GetAllPosts(ctx context.Context, filter post.PostFilter, page, limit int) (*post.PostsListResponse, error)
	GetMyPosts(ctx context.Context, authorID uint, page, limit int) (*post.PostsListResponse, error)
	GetPublishedPosts(ctx context.Context, page, limit int) (*post.PostsListResponse, error)
	PublishPost(ctx context.Context, id uint, userID uint, userRole string) (*post.PostResponse, error)
	UnpublishPost(ctx context.Context, id uint, userID uint, userRole string) (*post.PostResponse, error)
}

type postUseCase struct {
	postRepo post.Repository
	userRepo user.Repository
}

// NewPostUseCase creates a new post use case
func NewPostUseCase(postRepo post.Repository, userRepo user.Repository) PostUseCase {
	return &postUseCase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (uc *postUseCase) CreatePost(ctx context.Context, req post.CreatePostRequest, authorID uint) (*post.PostResponse, error) {
	// Generate slug from title
	slug := uc.generateSlug(req.Title)

	// Check if slug already exists
	existingPost, _ := uc.postRepo.GetBySlug(ctx, slug)
	if existingPost != nil {
		slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
	}

	// Set default values
	status := "draft"
	if req.Status != nil {
		status = *req.Status
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	// Create post
	newPost := &post.Post{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Slug:        slug,
		Status:      status,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
		FeaturedImg: req.FeaturedImg,
		IsPublic:    isPublic,
	}

	// Set published_at if status is published
	if status == "published" {
		now := time.Now()
		newPost.PublishedAt = &now
	}

	if err := uc.postRepo.Create(ctx, newPost); err != nil {
		return nil, errors.New("failed to create post")
	}

	return uc.mapToPostResponse(ctx, newPost)
}

func (uc *postUseCase) GetPostByID(ctx context.Context, id uint, incrementView bool) (*post.PostResponse, error) {
	p, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Increment view count if requested
	if incrementView {
		uc.postRepo.IncrementViewCount(ctx, id)
		p.ViewCount++
	}

	return uc.mapToPostResponse(ctx, p)
}

func (uc *postUseCase) GetPostBySlug(ctx context.Context, slug string, incrementView bool) (*post.PostResponse, error) {
	p, err := uc.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Increment view count if requested
	if incrementView {
		uc.postRepo.IncrementViewCount(ctx, p.ID)
		p.ViewCount++
	}

	return uc.mapToPostResponse(ctx, p)
}

func (uc *postUseCase) UpdatePost(ctx context.Context, id uint, req post.UpdatePostRequest, userID uint, userRole string) (*post.PostResponse, error) {
	p, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !uc.canModifyPost(p, userID, userRole) {
		return nil, errors.New("permission denied")
	}

	// Update fields if provided
	if req.Title != nil {
		p.Title = *req.Title
		// Regenerate slug if title changed
		newSlug := uc.generateSlug(*req.Title)
		if newSlug != p.Slug {
			// Check if new slug exists
			existingPost, _ := uc.postRepo.GetBySlug(ctx, newSlug)
			if existingPost != nil && existingPost.ID != p.ID {
				newSlug = fmt.Sprintf("%s-%d", newSlug, time.Now().Unix())
			}
			p.Slug = newSlug
		}
	}

	if req.Content != nil {
		p.Content = *req.Content
	}

	if req.Summary != nil {
		p.Summary = req.Summary
	}

	if req.CategoryID != nil {
		p.CategoryID = req.CategoryID
	}

	if req.FeaturedImg != nil {
		p.FeaturedImg = req.FeaturedImg
	}

	if req.IsPublic != nil {
		p.IsPublic = *req.IsPublic
	}

	if req.Status != nil {
		oldStatus := p.Status
		p.Status = *req.Status

		// Set published_at when changing to published
		if oldStatus != "published" && *req.Status == "published" {
			now := time.Now()
			p.PublishedAt = &now
		}
	}

	if err := uc.postRepo.Update(ctx, p); err != nil {
		return nil, errors.New("failed to update post")
	}

	return uc.mapToPostResponse(ctx, p)
}

func (uc *postUseCase) DeletePost(ctx context.Context, id uint, userID uint, userRole string) error {
	p, err := uc.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check permissions
	if !uc.canModifyPost(p, userID, userRole) {
		return errors.New("permission denied")
	}

	if err := uc.postRepo.Delete(ctx, id); err != nil {
		return errors.New("failed to delete post")
	}

	return nil
}

func (uc *postUseCase) GetAllPosts(ctx context.Context, filter post.PostFilter, page, limit int) (*post.PostsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	posts, err := uc.postRepo.GetAll(ctx, filter, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch posts")
	}

	total, err := uc.postRepo.GetTotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to count posts")
	}

	postResponses := make([]post.PostResponse, len(posts))
	for i, p := range posts {
		response, err := uc.mapToPostResponse(ctx, p)
		if err != nil {
			continue // Skip if error mapping
		}
		postResponses[i] = *response
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &post.PostsListResponse{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *postUseCase) GetMyPosts(ctx context.Context, authorID uint, page, limit int) (*post.PostsListResponse, error) {
	filter := post.PostFilter{
		AuthorID: &authorID,
	}
	return uc.GetAllPosts(ctx, filter, page, limit)
}

func (uc *postUseCase) GetPublishedPosts(ctx context.Context, page, limit int) (*post.PostsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	posts, err := uc.postRepo.GetPublished(ctx, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch published posts")
	}

	// Get total count for published posts
	publishedStatus := "published"
	isPublic := true
	filter := post.PostFilter{
		Status:   &publishedStatus,
		IsPublic: &isPublic,
	}
	total, err := uc.postRepo.GetTotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New("failed to count published posts")
	}

	postResponses := make([]post.PostResponse, len(posts))
	for i, p := range posts {
		response, err := uc.mapToPostResponse(ctx, p)
		if err != nil {
			continue
		}
		postResponses[i] = *response
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &post.PostsListResponse{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *postUseCase) PublishPost(ctx context.Context, id uint, userID uint, userRole string) (*post.PostResponse, error) {
	req := post.UpdatePostRequest{
		Status: stringPtr("published"),
	}
	return uc.UpdatePost(ctx, id, req, userID, userRole)
}

func (uc *postUseCase) UnpublishPost(ctx context.Context, id uint, userID uint, userRole string) (*post.PostResponse, error) {
	req := post.UpdatePostRequest{
		Status: stringPtr("draft"),
	}
	return uc.UpdatePost(ctx, id, req, userID, userRole)
}

// Helper functions
func (uc *postUseCase) generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

func (uc *postUseCase) canModifyPost(p *post.Post, userID uint, userRole string) bool {
	// Admin can modify any post
	if userRole == "admin" {
		return true
	}

	// Author can modify their own post
	return p.AuthorID == userID
}

func (uc *postUseCase) mapToPostResponse(ctx context.Context, p *post.Post) (*post.PostResponse, error) {
	// Get author name
	author, err := uc.userRepo.GetByID(ctx, p.AuthorID)
	authorName := "Unknown"
	if err == nil && author != nil {
		authorName = author.Name
	}

	// Handle nil pointers
	summary := ""
	if p.Summary != nil {
		summary = *p.Summary
	}

	featuredImg := ""
	if p.FeaturedImg != nil {
		featuredImg = *p.FeaturedImg
	}

	return &post.PostResponse{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Summary:     summary,
		Slug:        p.Slug,
		Status:      p.Status,
		CategoryID:  p.CategoryID,
		AuthorID:    p.AuthorID,
		AuthorName:  authorName,
		FeaturedImg: featuredImg,
		ViewCount:   p.ViewCount,
		IsPublic:    p.IsPublic,
		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func stringPtr(s string) *string {
	return &s
}
