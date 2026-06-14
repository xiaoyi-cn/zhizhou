package parser

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
)

// ParseWebContent 提取网页正文
func ParseWebContent(rawURL string) (string, string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ZhiZhou/1.0; +https://zhizhou.app)")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := article.RenderText(&buf); err != nil {
		return article.Title(), "", nil
	}

	return article.Title(), buf.String(), nil
}