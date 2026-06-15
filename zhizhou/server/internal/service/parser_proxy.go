package service

import "github.com/your-username/zhizhou/server/internal/pkg/parser"

// ParseWebContent 代理 parser 包的函数，方便依赖注入
func ParseWebContent(rawURL string) (string, string, error) {
	return parser.ParseWebContent(rawURL)
}