package repository

import (
	"context"
	"errors"
	"strings"

	"moon/internal/domain/post"

	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) post.Repository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Create(ctx context.Context, p *post.Post) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *postRepository) GetByID(ctx context.Context, id uint) (*post.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postRepository) GetBySlug(ctx context.Context, slug string) (*post.Post, error) {
	var p post.Post
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postRepository) Update(ctx context.Context, p *post.Post) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&post.Post{}, id).Error
}

func (r *postRepository) GetAll(ctx context.Context, filter post.PostFilter, limit, offset int) ([]*post.Post, error) {
	var posts []*post.Post
	query := r.db.WithContext(ctx).Model(&post.Post{})

	// Apply filters
	query = r.applyFilters(query, filter)

	err := query.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

func (r *postRepository) GetTotalCount(ctx context.Context, filter post.PostFilter) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&post.Post{})

	// Apply filters
	query = r.applyFilters(query, filter)

	err := query.Count(&count).Error
	return count, err
}

func (r *postRepository) GetByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*post.Post, error) {
	var posts []*post.Post
	err := r.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) GetByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]*post.Post, error) {
	var posts []*post.Post
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) GetPublished(ctx context.Context, limit, offset int) ([]*post.Post, error) {
	var posts []*post.Post
	err := r.db.WithContext(ctx).
		Where("status = ? AND is_public = ?", "published", true).
		Limit(limit).
		Offset(offset).
		Order("published_at DESC").
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&post.Post{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// Helper function to apply filters
func (r *postRepository) applyFilters(query *gorm.DB, filter post.PostFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}

	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?", searchTerm, searchTerm)
	}

	return query
}
