package model

import "time"

// User 用户
type User struct {
	ID           string    `json:"id"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Tier         string    `json:"tier"`
	StorageMode  string    `json:"storage_mode"`
	StorageUsed  int64     `json:"storage_used"`
	ProExpiresAt *time.Time `json:"pro_expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// APIKey 用户 API Key 配置
type APIKey struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Provider        string    `json:"provider"`
	APIKeyEncrypted string    `json:"-"`
	BaseURL         string    `json:"base_url"`
	Model           string    `json:"model"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// Content 内容
type Content struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	SourceType string    `json:"source_type"`
	RawContent string    `json:"raw_content,omitempty"`
	Summary    string    `json:"summary"`
	Category   string    `json:"category"`
	Tags       []string  `json:"tags"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Category 分类
type Category struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Subscription 订阅
type Subscription struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Tier        string     `json:"tier"`
	Plan        string     `json:"plan"`
	Status      string     `json:"status"`
	StartedAt   time.Time  `json:"started_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CancelledAt *time.Time `json:"cancelled_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// FeatureInfo 功能信息
type FeatureInfo struct {
	Tier     string          `json:"tier"`
	Features map[string]bool `json:"features"`
}

// SubscriptionInfo 订阅信息
type SubscriptionInfo struct {
	Tier         string          `json:"tier"`
	StorageUsed  int64           `json:"storage_used"`
	StorageLimit int64           `json:"storage_limit"`
	ProExpiresAt *time.Time      `json:"pro_expires_at"`
	Features     map[string]bool `json:"features"`
}