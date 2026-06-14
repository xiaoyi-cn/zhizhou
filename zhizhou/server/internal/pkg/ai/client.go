package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client LLM 客户端接口
type Client interface {
	Chat(systemPrompt, userPrompt string) (string, error)
	Summarize(title, content string) (string, error)
	Classify(categories []string, title, summary string) (string, error)
	ExtractTags(title, summary string) ([]string, error)
	GenerateEmbedding(text string) ([]float64, error)
}

// openaiClient OpenAI 兼容客户端
type openaiClient struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// ChatMessage chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest chat request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

// ChatResponse chat response
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// EmbeddingRequest embedding request
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbeddingResponse embedding response
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// NewClient 创建 OpenAI 兼容的 LLM 客户端
func NewClient(apiKey, baseURL, model string) Client {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &openaiClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *openaiClient) Chat(systemPrompt, userPrompt string) (string, error) {
	req := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API error: %s, body: %s", resp.Status, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *openaiClient) Summarize(title, content string) (string, error) {
	systemPrompt := "你是一个知识管理助手。请用 200 字以内总结以下文章的核心内容。"
	userPrompt := fmt.Sprintf("标题：%s\n正文：%s", title, content)
	return c.Chat(systemPrompt, userPrompt)
}

func (c *openaiClient) Classify(categories []string, title, summary string) (string, error) {
	systemPrompt := "你是一个知识管理助手。请判断文章属于哪个分类。"
	userPrompt := fmt.Sprintf("用户已有的分类：%v\n请判断以下文章属于哪个分类。如果匹配已有分类，返回分类名。如果不匹配，建议一个新分类名（格式如 技术/后端）。\n标题：%s\n摘要：%s", categories, title, summary)
	return c.Chat(systemPrompt, userPrompt)
}

func (c *openaiClient) ExtractTags(title, summary string) ([]string, error) {
	systemPrompt := "你是一个知识管理助手。从文章提取标签。"
	userPrompt := fmt.Sprintf("从以下文章提取 3-5 个标签。标签为简短关键词如 Go K8s 微服务 开源。返回 JSON 数组格式。\n标题：%s\n摘要：%s", title, summary)
	result, err := c.Chat(systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}
	var tags []string
	if err := json.Unmarshal([]byte(result), &tags); err != nil {
		return []string{result}, nil
	}
	return tags, nil
}

func (c *openaiClient) GenerateEmbedding(text string) ([]float64, error) {
	req := EmbeddingRequest{
		Model: "text-embedding-ada-002",
		Input: text,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Embedding API error: %s", resp.Status)
	}

	var embResp EmbeddingResponse
	if err := json.Unmarshal(respBody, &embResp); err != nil {
		return nil, err
	}

	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embResp.Data[0].Embedding, nil
}