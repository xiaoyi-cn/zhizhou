package service

import (
	"github.com/your-username/zhizhou/server/internal/model"
	"github.com/your-username/zhizhou/server/internal/pkg/ai"
	"github.com/your-username/zhizhou/server/internal/pkg/crypto"
	"github.com/your-username/zhizhou/server/internal/repository"
)

// AuthService 认证服务
type AuthService interface {
	SendCode(phone string) error
	VerifyCode(phone, code string) (string, error) // returns userID
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) SendCode(phone string) error {
	// TODO: 接入阿里云短信/腾讯云短信
	return nil
}

func (s *authService) VerifyCode(phone, code string) (string, error) {
	// TODO: 验证短信验证码
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", err
	}
	if user == nil {
		// 新用户自动注册
		user = &model.User{Phone: phone}
		if err := s.userRepo.Create(user); err != nil {
			return "", err
		}
	}
	return user.ID, nil
}

// APIKeyService API Key 服务
type APIKeyService interface {
	SaveAPIKey(userID, provider, plainKey, baseURL, modelName string) error
	ListAPIKeys(userID string) ([]model.APIKey, error)
	GetAPIKey(id string) (*model.APIKey, error)
	UpdateAPIKey(id, provider, plainKey, baseURL, modelName string) error
	DeleteAPIKey(id string) error
	TestAPIKey(provider, plainKey, baseURL, modelName string) error
	GetDecryptedKey(userID string) (string, string, string, error) // APIKey, BaseURL, Model
}

type apiKeyService struct {
	repo    repository.APIKeyRepository
	encKey  []byte
}

func NewAPIKeyService(repo repository.APIKeyRepository, encKey string) APIKeyService {
	return &apiKeyService{
		repo:   repo,
		encKey: []byte(encKey),
	}
}

func (s *apiKeyService) SaveAPIKey(userID, provider, plainKey, baseURL, modelName string) error {
	encrypted, err := crypto.Encrypt(s.encKey, plainKey)
	if err != nil {
		return err
	}
	k := &model.APIKey{
		UserID:          userID,
		Provider:        provider,
		APIKeyEncrypted: encrypted,
		BaseURL:         baseURL,
		Model:           modelName,
	}
	return s.repo.Create(k)
}

func (s *apiKeyService) ListAPIKeys(userID string) ([]model.APIKey, error) {
	return s.repo.GetByUserID(userID)
}

func (s *apiKeyService) GetAPIKey(id string) (*model.APIKey, error) {
	return s.repo.GetByID(id)
}

func (s *apiKeyService) UpdateAPIKey(id, provider, plainKey, baseURL, modelName string) error {
	encrypted, err := crypto.Encrypt(s.encKey, plainKey)
	if err != nil {
		return err
	}
	k := &model.APIKey{
		ID:              id,
		Provider:        provider,
		APIKeyEncrypted: encrypted,
		BaseURL:         baseURL,
		Model:           modelName,
		IsActive:        true,
	}
	return s.repo.Update(k)
}

func (s *apiKeyService) DeleteAPIKey(id string) error {
	return s.repo.Delete(id)
}

func (s *apiKeyService) TestAPIKey(provider, plainKey, baseURL, modelName string) error {
	client := ai.NewClient(plainKey, baseURL, modelName)
	_, err := client.Chat("Hello", "Say hi")
	return err
}

func (s *apiKeyService) GetDecryptedKey(userID string) (string, string, string, error) {
	keys, err := s.repo.GetByUserID(userID)
	if err != nil || len(keys) == 0 {
		return "", "", "", err
	}
	// 取第一个激活的 key
	var activeKey *model.APIKey
	for _, k := range keys {
		if k.IsActive {
			activeKey = &k
			break
		}
	}
	if activeKey == nil {
		activeKey = &keys[0]
	}
	decrypted, err := crypto.Decrypt(s.encKey, activeKey.APIKeyEncrypted)
	if err != nil {
		return "", "", "", err
	}
	return decrypted, activeKey.BaseURL, activeKey.Model, nil
}

// IngestService 内容采集服务
type IngestService interface {
	ProcessURL(userID, rawURL string) (*model.Content, error)
}

type ingestService struct {
	contentRepo repository.ContentRepository
	apiKeySvc   APIKeyService
}

func NewIngestService(contentRepo repository.ContentRepository, apiKeySvc APIKeyService) IngestService {
	return &ingestService{
		contentRepo: contentRepo,
		apiKeySvc:   apiKeySvc,
	}
}

func (s *ingestService) ProcessURL(userID, rawURL string) (*model.Content, error) {
	// TODO: Phase 2 改为异步任务队列
	// 1. 抓取正文
	title, content, err := ParseWebContent(rawURL)
	if err != nil {
		// 降级：解析失败也入库，标题用 URL
		title = rawURL
		content = ""
	}

	// 2. 入库（status=pending）
	item := &model.Content{
		UserID:     userID,
		URL:        rawURL,
		Title:      title,
		SourceType: "web",
		RawContent: content,
		Status:     "pending",
	}
	if err := s.contentRepo.Create(item); err != nil {
		return nil, err
	}

	// 3. 异步 AI 处理（简化：在入库后直接调用）
	go s.processWithAI(item)

	return item, nil
}

func (s *ingestService) processWithAI(c *model.Content) {
	apiKey, baseURL, model, err := s.apiKeySvc.GetDecryptedKey(c.UserID)
	if err != nil || apiKey == "" {
		return
	}

	client := ai.NewClient(apiKey, baseURL, model)

	// 摘要
	summary, err := client.Summarize(c.Title, c.RawContent)
	if err == nil {
		c.Summary = summary
	}

	// 标签
	tags, err := client.ExtractTags(c.Title, c.Summary)
	if err == nil {
		c.Tags = tags
	}

	// 分类
	category, err := client.Classify(nil, c.Title, c.Summary)
	if err == nil {
		c.Category = category
	}

	s.contentRepo.Update(c)
}

// ContentService 内容服务
type ContentService interface {
	GetPending(userID string) ([]model.Content, error)
	Approve(id string) error
	Skip(id string) error
	Update(id string, title, summary, category string, tags []string) error
	GetByID(id string) (*model.Content, error)
	List(userID string, tags []string, category string, page, limit int) ([]model.Content, int, error)
	Search(userID, query string, tags []string, page, limit int) ([]model.Content, int, error)
}

type contentService struct {
	repo repository.ContentRepository
}

func NewContentService(repo repository.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

func (s *contentService) GetPending(userID string) ([]model.Content, error) {
	return s.repo.GetPendingByUserID(userID)
}

func (s *contentService) Approve(id string) error {
	return s.repo.Approve(id)
}

func (s *contentService) Skip(id string) error {
	return s.repo.Skip(id)
}

func (s *contentService) Update(id string, title, summary, category string, tags []string) error {
	c := &model.Content{
		ID:       id,
		Title:    title,
		Summary:  summary,
		Category: category,
		Tags:     tags,
	}
	return s.repo.Update(c)
}

func (s *contentService) GetByID(id string) (*model.Content, error) {
	return s.repo.GetByID(id)
}

func (s *contentService) List(userID string, tags []string, category string, page, limit int) ([]model.Content, int, error) {
	return s.repo.List(userID, tags, category, page, limit)
}

func (s *contentService) Search(userID, query string, tags []string, page, limit int) ([]model.Content, int, error) {
	return s.repo.Search(userID, query, tags, page, limit)
}

// CategoryService 分类服务
type CategoryService interface {
	Create(userID, name string, parentID *string) (*model.Category, error)
	List(userID string) ([]model.Category, error)
	Delete(id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(userID, name string, parentID *string) (*model.Category, error) {
	c := &model.Category{
		UserID:   userID,
		Name:     name,
		ParentID: parentID,
	}
	return c, s.repo.Create(c)
}

func (s *categoryService) List(userID string) ([]model.Category, error) {
	return s.repo.ListByUserID(userID)
}

func (s *categoryService) Delete(id string) error {
	return s.repo.Delete(id)
}

// SubscriptionService 订阅服务
type SubscriptionService interface {
	GetSubscriptionInfo(userID, deployMode string) (*model.SubscriptionInfo, error)
	GetFeatures(deployMode, tier string) *model.FeatureInfo
}

type subscriptionService struct {
	userRepo repository.UserRepository
}

func NewSubscriptionService(userRepo repository.UserRepository) SubscriptionService {
	return &subscriptionService{userRepo: userRepo}
}

func (s *subscriptionService) GetSubscriptionInfo(userID, deployMode string) (*model.SubscriptionInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	tier := user.Tier
	if deployMode == "local" {
		tier = "pro"
	}

	info := &model.SubscriptionInfo{
		Tier:         tier,
		StorageUsed:  user.StorageUsed,
		StorageLimit: 10737418240, // 10GB for free
		ProExpiresAt: user.ProExpiresAt,
		Features:     s.GetFeatures(deployMode, tier).Features,
	}

	if tier == "pro" {
		info.StorageLimit = 0 // unlimited
	}

	return info, nil
}

func (s *subscriptionService) GetFeatures(deployMode, tier string) *model.FeatureInfo {
	features := map[string]bool{
		"mcp_server": false,
		"open_api":   false,
		"auto_rules": false,
	}

	if deployMode == "local" || tier == "pro" {
		features["mcp_server"] = true
		features["open_api"] = true
		features["auto_rules"] = true
	}

	return &model.FeatureInfo{
		Tier:     tier,
		Features: features,
	}
}