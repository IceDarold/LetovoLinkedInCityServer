package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// NotificationService отправляет сообщения в несколько Telegram чатов
type NotificationService struct {
	apiURL  string
	chatIDs []int64
}

// NewNotificationService создаёт NotificationService
func NewNotificationService() *NotificationService {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		panic("NotificationService: TELEGRAM_BOT_TOKEN not set")
	}
	chatIDsStr := os.Getenv("TELEGRAM_CHAT_IDS")
	if chatIDsStr == "" {
		panic("NotificationService: TELEGRAM_CHAT_IDS not set")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	chatIDs, err := parseChatIDs(chatIDsStr)
	if err != nil {
		panic("NotificationService: invalid TELEGRAM_CHAT_IDS: " + err.Error())
	}

	return &NotificationService{
		apiURL:  apiURL,
		chatIDs: chatIDs,
	}
}

// вспомогательная функция
func parseChatIDs(raw string) ([]int64, error) {
	parts := strings.Split(raw, ",")
	var result []int64
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var id int64
		_, err := fmt.Sscan(part, &id)
		if err != nil {
			return nil, fmt.Errorf("не удалось разобрать chat_id: %s", part)
		}
		result = append(result, id)
	}
	if len(result) == 0 {
		return nil, errors.New("список chat_ids пустой")
	}
	return result, nil
}

type telegramRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

type telegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	Description string          `json:"description,omitempty"`
}

// SendMessage отправляет текст во все чаты
func (ns *NotificationService) SendMessage(text string) error {
	var lastErr error

	for _, chatID := range ns.chatIDs {
		reqBody := telegramRequest{
			ChatID:    chatID,
			Text:      text,
			ParseMode: "HTML",
		}
		payload, err := json.Marshal(reqBody)
		if err != nil {
			lastErr = fmt.Errorf("marshal telegram request: %w", err)
			continue
		}

		resp, err := http.Post(ns.apiURL, "application/json", bytes.NewReader(payload))
		if err != nil {
			lastErr = fmt.Errorf("post to telegram API: %w", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var tgResp telegramResponse
		if err := json.Unmarshal(body, &tgResp); err != nil {
			lastErr = fmt.Errorf("unmarshal telegram response: %w", err)
			continue
		}
		if !tgResp.Ok {
			desc := tgResp.Description
			if desc == "" {
				desc = fmt.Sprintf("HTTP status %d", resp.StatusCode)
			}
			lastErr = errors.New("telegram API error: " + desc)
		}
	}

	return lastErr
}
