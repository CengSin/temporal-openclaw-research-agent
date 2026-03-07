package activity

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a3tai/openclaw-go/chatcompletions"
)

// GenerateArticleWithChronicle 调用本地 OpenClaw 中的 Chronicle agent 生成简短文章。
func GenerateArticleWithChronicle(ctx context.Context, topic string) (string, error) {
	baseURL := os.Getenv("OPENCLAW_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:18789"
	}

	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	client := &chatcompletions.Client{
		BaseURL:    baseURL,
		Token:      os.Getenv("OPENCLAW_TOKEN"),
		AgentID:    "Chronicle",
		SessionKey: "temporal-chronicle",
	}

	prompt := fmt.Sprintf(
		"请围绕主题《%s》调用可用 skills，写一篇简短中文文章。输出格式：标题一行 + 3 段正文，每段不超过 80 字。",
		topic,
	)

	resp, err := client.Create(reqCtx, chatcompletions.Request{
		Model: "openclaw:chronicle",
		Messages: []chatcompletions.Message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		if httpErr, ok := errors.AsType[*chatcompletions.HTTPError](err); ok {
			switch httpErr.StatusCode {
			case 401:
				return "", fmt.Errorf(
					"调用 Chronicle 失败: OpenClaw 需要鉴权，请设置 OPENCLAW_TOKEN（baseURL=%s）: %w",
					baseURL,
					err,
				)
			case 404:
				return "", fmt.Errorf(
					"调用 Chronicle 失败: %s 没有暴露 /v1/chat/completions，当前地址更像 OpenClaw Control UI，请检查 OPENCLAW_BASE_URL 或本地 API 能力: %w",
					baseURL,
					err,
				)
			}
		}
		return "", fmt.Errorf("调用 Chronicle 失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("Chronicle 返回为空（choices=0）")
	}

	article := strings.TrimSpace(resp.Choices[0].Message.Content)
	if article == "" {
		return "", errors.New("Chronicle 返回内容为空")
	}

	return article, nil
}
