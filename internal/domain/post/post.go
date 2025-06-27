package post

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content" gorm:"type:text"`
	Summary     *string        `json:"summary" gorm:"type:text"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Status      string         `json:"status" gorm:"default:'draft'"` // draft, published, archived
	CategoryID  *uint          `json:"category_id"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	FeaturedImg *string        `json:"featured_img"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePostRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=200"`
	Content     string  `json:"content" binding:"required"`
	Summary     *string `json:"summary"`
	CategoryID  *uint   `json:"category_id"`
	FeaturedImg *string `json:"featured_img"`
	IsPublic    *bool   `json:"is_public"`
	Status      *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
	Content     *string `json:"content"`
	Summary     *string `json:"summary"`
	CategoryID  *uint   `json:"category_id"`
	FeaturedImg *string `json:"featured_img"`
	IsPublic    *bool   `json:"is_public"`
	Status      *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type PostResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Summary     string     `json:"summary"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	CategoryID  *uint      `json:"category_id"`
	AuthorID    uint       `json:"author_id"`
	AuthorName  string     `json:"author_name"`
	FeaturedImg string     `json:"featured_img"`
	ViewCount   int        `json:"view_count"`
	IsPublic    bool       `json:"is_public"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostsListResponse struct {
	Posts      []PostResponse `json:"posts"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type PostFilter struct {
	Status     *string `json:"status"`
	CategoryID *uint   `json:"category_id"`
	AuthorID   *uint   `json:"author_id"`
	IsPublic   *bool   `json:"is_public"`
	Search     *string `json:"search"` // Search in title and content
}

// Repository interface - Domain layer
type Repository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	GetBySlug(ctx context.Context, slug string) (*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context, filter PostFilter, limit, offset int) ([]*Post, error)
	GetTotalCount(ctx context.Context, filter PostFilter) (int64, error)
	GetByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]*Post, error)
	GetByCategory(ctx context.Context, categoryID uint, limit, offset int) ([]*Post, error)
	GetPublished(ctx context.Context, limit, offset int) ([]*Post, error)
	IncrementViewCount(ctx context.Context, id uint) error
}
