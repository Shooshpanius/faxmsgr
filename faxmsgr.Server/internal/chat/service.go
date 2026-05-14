// Пакет chat управляет чатами и сообщениями.
package chat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Chat представляет чат (диалог или группу).
type Chat struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"` // "direct" или "group"
	Name          string    `json:"name,omitempty"`
	LastMessage   string    `json:"last_message,omitempty"`
	LastMessageAt time.Time `json:"last_message_at,omitempty"`
}

// CreateChatReq параметры создания чата.
type CreateChatReq struct {
	Type      string   `json:"type"`       // "direct" или "group"
	Name      string   `json:"name"`       // для группового чата
	MemberIDs []string `json:"member_ids"` // участники
}

// Service управляет чатами.
type Service struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewService создаёт сервис чатов.
func NewService(db *pgxpool.Pool, logger *zap.Logger) *Service {
	return &Service{db: db, logger: logger}
}

// CreateChat создаёт новый чат и добавляет участников.
func (s *Service) CreateChat(ctx context.Context, creatorID string, req CreateChatReq) (*Chat, error) {
	if req.Type != "direct" && req.Type != "group" {
		return nil, errors.New("type must be 'direct' or 'group'")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var chatID string
	err = tx.QueryRow(ctx,
		`INSERT INTO chats (type, name) VALUES ($1, $2) RETURNING id`,
		req.Type, req.Name,
	).Scan(&chatID)
	if err != nil {
		return nil, fmt.Errorf("insert chat: %w", err)
	}

	// Добавить создателя в участники
	members := append([]string{creatorID}, req.MemberIDs...)
	for _, uid := range members {
		if _, err := tx.Exec(ctx,
			`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			chatID, uid,
		); err != nil {
			return nil, fmt.Errorf("insert member: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &Chat{ID: chatID, Type: req.Type, Name: req.Name}, nil
}

// ListChats возвращает список активных чатов пользователя с последним сообщением.
func (s *Service) ListChats(ctx context.Context, userID string) ([]*Chat, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.type, COALESCE(c.name, ''),
		        COALESCE(m.body, ''), COALESCE(m.created_at, '1970-01-01')
		 FROM chats c
		 JOIN chat_members cm ON cm.chat_id = c.id AND cm.user_id = $1
		 LEFT JOIN LATERAL (
		     SELECT body, created_at FROM messages
		     WHERE chat_id = c.id ORDER BY created_at DESC LIMIT 1
		 ) m ON true
		 ORDER BY COALESCE(m.created_at, '1970-01-01') DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list chats: %w", err)
	}
	defer rows.Close()

	var chats []*Chat
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID, &c.Type, &c.Name, &c.LastMessage, &c.LastMessageAt); err != nil {
			return nil, err
		}
		chats = append(chats, &c)
	}
	return chats, rows.Err()
}
