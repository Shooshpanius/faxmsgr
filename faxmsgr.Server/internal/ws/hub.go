// Пакет ws реализует WebSocket-хаб для обмена сообщениями в реальном времени.
package ws

import (
	"sync"

	"go.uber.org/zap"
)

// Hub управляет активными WebSocket-соединениями.
type Hub struct {
	// clients хранит все подключённые клиенты (userID -> *Client)
	clients map[string]*Client
	mu      sync.RWMutex

	// broadcast канал для рассылки сообщений
	broadcast chan *Message

	// register/unregister регистрация/снятие клиентов
	register   chan *Client
	unregister chan *Client

	logger *zap.Logger
}

// Message пакет WebSocket-сообщения.
type Message struct {
	Type    string `json:"type"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	ChatID  string `json:"chat_id,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// NewHub создаёт новый хаб.
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client, 16),
		unregister: make(chan *Client, 16),
		logger:     logger,
	}
}

// Run запускает основной цикл обработки событий хаба.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			h.logger.Info("ws client connected", zap.String("userID", client.userID))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Info("ws client disconnected", zap.String("userID", client.userID))

		case msg := <-h.broadcast:
			h.mu.RLock()
			target, online := h.clients[msg.To]
			h.mu.RUnlock()
			if online {
				select {
				case target.send <- msg:
				default:
					// Буфер переполнен — отключить клиента
					h.mu.Lock()
					delete(h.clients, target.userID)
					close(target.send)
					h.mu.Unlock()
				}
			}
		}
	}
}

// IsOnline проверяет, подключён ли пользователь.
func (h *Hub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// Send отправляет сообщение конкретному пользователю.
func (h *Hub) Send(msg *Message) {
	h.broadcast <- msg
}
