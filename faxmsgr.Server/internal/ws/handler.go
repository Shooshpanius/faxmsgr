package ws

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// В продакшене следует проверять Origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

// MakeHandler возвращает HTTP-обработчик для WebSocket-соединений.
// JWT-токен передаётся в query-параметре ?token=...
func MakeHandler(hub *Hub, jwtSecret string, logger *zap.Logger) http.HandlerFunc {
	secret := []byte(jwtSecret)
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		userID, err := extractUserID(tokenStr, secret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("ws upgrade failed", zap.Error(err))
			return
		}

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan *Message, 256),
			userID: userID,
			logger: logger,
		}
		hub.register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}

// extractUserID извлекает userID из JWT-токена.
func extractUserID(tokenStr string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("missing sub")
	}
	return sub, nil
}
