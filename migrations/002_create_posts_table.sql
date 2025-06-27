-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    status ENUM('draft', 'published', 'archived') DEFAULT 'draft',
    category_id INT NULL,
    author_id INT NOT NULL,
    featured_img VARCHAR(500) NULL,
    view_count INT DEFAULT 0,
    is_public BOOLEAN DEFAULT TRUE,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Foreign key constraint (if users table exists)
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes for better performance
    INDEX idx_posts_slug (slug),
    INDEX idx_posts_status (status),
    INDEX idx_posts_author_id (author_id),
    INDEX idx_posts_category_id (category_id),
    INDEX idx_posts_published_at (published_at),
    INDEX idx_posts_deleted_at (deleted_at),
    INDEX idx_posts_status_public (status, is_public)
);

-- Create full-text index for search functionality
ALTER TABLE posts ADD FULLTEXT(title, content); 