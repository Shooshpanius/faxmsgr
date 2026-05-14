package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

// Client представляет одно WebSocket-соединение.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *Message
	userID string
	logger *zap.Logger
}

// ReadPump читает входящие сообщения от клиента.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) //nolint:errcheck
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Warn("ws read error", zap.String("userID", c.userID), zap.Error(err))
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			c.logger.Warn("ws invalid message", zap.String("userID", c.userID), zap.Error(err))
			continue
		}
		msg.From = c.userID

		// Обработка типов сообщений
		switch msg.Type {
		case "MSG_SEND":
			// Перенаправить получателю (если онлайн)
			c.hub.Send(&msg)
		case "STATUS_ONLINE", "STATUS_OFFLINE":
			// Статус обрабатывается на уровне подключения/отключения
			c.logger.Info("ws status event", zap.String("type", msg.Type), zap.String("userID", c.userID))
		}
	}
}

// WritePump отправляет исходящие сообщения клиенту.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //nolint:errcheck
				return
			}
			data, _ := json.Marshal(msg)
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //nolint:errcheck
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
