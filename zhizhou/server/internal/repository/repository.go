package repository

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/your-username/zhizhou/server/internal/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *model.User) error
	FindByPhone(phone string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	UpdateStorageUsed(userID string, delta int64) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (phone) VALUES ($1) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, user.Phone).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) FindByPhone(phone string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, phone, tier, storage_mode, storage_used, pro_expires_at, created_at, updated_at FROM users WHERE phone = $1`
	err := r.db.QueryRow(query, phone).Scan(
		&user.ID, &user.Phone, &user.Tier, &user.StorageMode,
		&user.StorageUsed, &user.ProExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepo) FindByID(id string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, phone, tier, storage_mode, storage_used, pro_expires_at, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Phone, &user.Tier, &user.StorageMode,
		&user.StorageUsed, &user.ProExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *userRepo) UpdateStorageUsed(userID string, delta int64) error {
	_, err := r.db.Exec(`UPDATE users SET storage_used = storage_used + $1, updated_at = NOW() WHERE id = $2`, delta, userID)
	return err
}

// APIKeyRepository API Key 仓储接口
type APIKeyRepository interface {
	Create(apiKey *model.APIKey) error
	GetByUserID(userID string) ([]model.APIKey, error)
	GetByID(id string) (*model.APIKey, error)
	Update(apiKey *model.APIKey) error
	Delete(id string) error
}

type apiKeyRepo struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) APIKeyRepository {
	return &apiKeyRepo{db: db}
}

func (r *apiKeyRepo) Create(k *model.APIKey) error {
	query := `INSERT INTO api_keys (user_id, provider, api_key_encrypted, base_url, model) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, k.UserID, k.Provider, k.APIKeyEncrypted, k.BaseURL, k.Model).Scan(&k.ID)
}

func (r *apiKeyRepo) GetByUserID(userID string) ([]model.APIKey, error) {
	rows, err := r.db.Query(`SELECT id, provider, api_key_encrypted, base_url, model, is_active 
	                        FROM api_keys WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []model.APIKey
	for rows.Next() {
		var k model.APIKey
		err := rows.Scan(&k.ID, &k.Provider, &k.APIKeyEncrypted, &k.BaseURL, &k.Model, &k.IsActive)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *apiKeyRepo) GetByID(id string) (*model.APIKey, error) {
	k := &model.APIKey{}
	query := `SELECT id, user_id, provider, api_key_encrypted, base_url, model, is_active FROM api_keys WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&k.ID, &k.UserID, &k.Provider, &k.APIKeyEncrypted, &k.BaseURL, &k.Model, &k.IsActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return k, err
}

func (r *apiKeyRepo) Update(k *model.APIKey) error {
	_, err := r.db.Exec(`UPDATE api_keys SET provider=$1, api_key_encrypted=$2, base_url=$3, model=$4, is_active=$5 WHERE id=$6`,
		k.Provider, k.APIKeyEncrypted, k.BaseURL, k.Model, k.IsActive, k.ID)
	return err
}

func (r *apiKeyRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM api_keys WHERE id = $1`, id)
	return err
}

// ContentRepository 内容仓储接口
type ContentRepository interface {
	Create(c *model.Content) error
	GetPendingByUserID(userID string) ([]model.Content, error)
	GetByID(id string) (*model.Content, error)
	Update(c *model.Content) error
	Approve(id string) error
	Skip(id string) error
	List(userID string, tags []string, category string, page, limit int) ([]model.Content, int, error)
	Search(userID, query string, tags []string, page, limit int) ([]model.Content, int, error)
}

type contentRepo struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) ContentRepository {
	return &contentRepo{db: db}
}

func (r *contentRepo) Create(c *model.Content) error {
	query := `
		INSERT INTO contents (user_id, url, title, source_type, raw_content, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		c.UserID, c.URL, c.Title, c.SourceType, c.RawContent, c.Status,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *contentRepo) GetPendingByUserID(userID string) ([]model.Content, error) {
	query := `
		SELECT id, url, title, summary, category, tags, status, created_at 
		FROM contents 
		WHERE user_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(&c.ID, &c.URL, &c.Title, &c.Summary, &c.Category,
			pq.Array(&c.Tags), &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (r *contentRepo) GetByID(id string) (*model.Content, error) {
	c := &model.Content{}
	query := `SELECT id, user_id, url, title, source_type, raw_content, summary, category, tags, status, created_at, updated_at FROM contents WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.UserID, &c.URL, &c.Title, &c.SourceType, &c.RawContent,
		&c.Summary, &c.Category, pq.Array(&c.Tags), &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *contentRepo) Update(c *model.Content) error {
	_, err := r.db.Exec(`UPDATE contents SET title=$1, summary=$2, category=$3, tags=$4, updated_at=NOW() WHERE id=$5`,
		c.Title, c.Summary, c.Category, pq.Array(c.Tags), c.ID)
	return err
}

func (r *contentRepo) Approve(id string) error {
	_, err := r.db.Exec(`UPDATE contents SET status='approved', updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *contentRepo) Skip(id string) error {
	_, err := r.db.Exec(`UPDATE contents SET status='skipped', updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *contentRepo) List(userID string, tags []string, category string, page, limit int) ([]model.Content, int, error) {
	// Count query
	var total int
	baseQuery := `FROM contents WHERE user_id = $1 AND status = 'approved'`
	args := []interface{}{userID}
	argIdx := 2

	if category != "" {
		baseQuery += ` AND category = $` + fmt.Sprintf("%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if len(tags) > 0 {
		baseQuery += ` AND tags @> $` + fmt.Sprintf("%d", argIdx)
		args = append(args, pq.Array(tags))
		argIdx++
	}

	err := r.db.QueryRow(`SELECT COUNT(*) `+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Select query
	offset := (page - 1) * limit
	selectQuery := `SELECT id, url, title, summary, category, tags, status, created_at ` + baseQuery +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contents []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(&c.ID, &c.URL, &c.Title, &c.Summary, &c.Category,
			pq.Array(&c.Tags), &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		contents = append(contents, c)
	}
	return contents, total, nil
}

func (r *contentRepo) Search(userID, query string, tags []string, page, limit int) ([]model.Content, int, error) {
	// Full text search on title and summary
	searchTerm := "%" + query + "%"
	baseQuery := `FROM contents WHERE user_id = $1 AND status = 'approved' AND (title ILIKE $2 OR summary ILIKE $2)`
	args := []interface{}{userID, searchTerm}
	argIdx := 3

	if len(tags) > 0 {
		baseQuery += ` AND tags @> $` + fmt.Sprintf("%d", argIdx)
		args = append(args, pq.Array(tags))
		argIdx++
	}

	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) `+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	selectQuery := `SELECT id, url, title, summary, category, tags, status, created_at ` + baseQuery +
		` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIdx) + ` OFFSET $` + fmt.Sprintf("%d", argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var contents []model.Content
	for rows.Next() {
		var c model.Content
		err := rows.Scan(&c.ID, &c.URL, &c.Title, &c.Summary, &c.Category,
			pq.Array(&c.Tags), &c.Status, &c.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		contents = append(contents, c)
	}
	return contents, total, nil
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(c *model.Category) error
	ListByUserID(userID string) ([]model.Category, error)
	Delete(id string) error
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(c *model.Category) error {
	query := `INSERT INTO categories (user_id, name, parent_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(query, c.UserID, c.Name, c.ParentID).Scan(&c.ID, &c.CreatedAt)
}

func (r *categoryRepo) ListByUserID(userID string) ([]model.Category, error) {
	rows, err := r.db.Query(`SELECT id, user_id, name, parent_id, created_at FROM categories WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.ParentID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}