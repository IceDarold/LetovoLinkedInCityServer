package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// NotificationService отправляет сообщения в Telegram‑группу через бот‑API.
type NotificationService struct {
	apiURL string
	chatID int64
}

// NewNotificationService создаёт NotificationService,
// загружая токен бота и ID чата из конфига.
func NewNotificationService() *NotificationService {
	token := viper.GetString("notification.telegram_bot_token")
	if token == "" {
		log.Fatal("NotificationService: missing config value 'notification.telegram_bot_token'")
	}
	chatID := viper.GetInt64("notification.telegram_chat_id")
	if chatID == 0 {
		log.Fatal("NotificationService: missing or invalid 'notification.telegram_chat_id'")
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	return &NotificationService{
		apiURL: apiURL,
		chatID: chatID,
	}
}

// telegramRequest — payload для API sendMessage
type telegramRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// telegramResponse — ответ Telegram API
type telegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	Description string          `json:"description,omitempty"`
}

// SendMessage отправляет текстовое сообщение в указанный чат.
// Текст отправляется в режиме HTML.
func (ns *NotificationService) SendMessage(text string) error {
	reqBody := telegramRequest{
		ChatID:    ns.chatID,
		Text:      text,
		ParseMode: "HTML",
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal telegram request: %w", err)
	}

	resp, err := http.Post(ns.apiURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("post to telegram API: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tgResp telegramResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		return fmt.Errorf("unmarshal telegram response: %w", err)
	}
	if !tgResp.Ok {
		desc := tgResp.Description
		if desc == "" {
			desc = fmt.Sprintf("HTTP status %d", resp.StatusCode)
		}
		return errors.New("telegram API error: " + desc)
	}

	return nil
}
