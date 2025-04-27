package services

import (
	"io"
	"os"
	"strings"
)

// MultiLogger отправляет лог одновременно в stdout и в Telegram с разметкой
type MultiLogger struct {
	stdout io.Writer
	notify *NotificationService
}

// NewMultiLogger создаёт мульти логгер
func NewMultiLogger(notify *NotificationService) *MultiLogger {
	return &MultiLogger{
		stdout: os.Stdout,
		notify: notify,
	}
}

// Write реализует интерфейс io.Writer
func (m *MultiLogger) Write(p []byte) (n int, err error) {
	// Сначала пишем в консоль
	n, err = m.stdout.Write(p)
	if err != nil {
		return n, err
	}

	// Подготовим сообщение для Telegram
	message := strings.TrimSpace(string(p))
	if message == "" {
		return n, nil
	}

	// Ограничим длину для Телеги
	const telegramLimit = 4000
	if len(message) > telegramLimit {
		message = message[:telegramLimit] + "..."
	}

	// Форматирование сообщения по ключевым словам
	var prefix, formattedMessage string

	lowerMsg := strings.ToLower(message)

	switch {
	case strings.Contains(lowerMsg, "error"), strings.Contains(lowerMsg, "ошибка"), strings.Contains(lowerMsg, "fail"), strings.Contains(lowerMsg, "panic"):
		prefix = "⚠️ <b>Ошибка:</b> "
		formattedMessage = prefix + "<code>" + message + "</code>"

	case strings.Contains(lowerMsg, "warn"), strings.Contains(lowerMsg, "warning"):
		prefix = "⚠️ <b>Предупреждение:</b> "
		formattedMessage = prefix + "<code>" + message + "</code>"

	default:
		prefix = "ℹ️ <b>Инфо:</b> "
		formattedMessage = prefix + "<code>" + message + "</code>"
	}

	// Отправляем в Telegram в отдельной горутине
	go func(msg string) {
		_ = m.notify.SendMessage(msg)
	}(formattedMessage)

	return n, nil
}
