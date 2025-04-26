package middleware

import (
	"city-server/internal/services"
	"city-server/internal/utils"
	"fmt"
	"html"
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorMiddleware ловит паники и ошибки, шлет вам уведомление
func ErrorMiddleware(ns *services.NotificationService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					errText := fmt.Sprintf("Panic: %v\n%s", rec, debug.Stack())
					// шлём уведомление
					if err := ns.SendMessage("🚨 <b>Panic</b> on " + r.URL.Path + "\n<code>" + html.EscapeString(errText) + "</code>"); err != nil {
						log.Println("Notification error:", err)
					}
					// стандартный ответ
					utils.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
