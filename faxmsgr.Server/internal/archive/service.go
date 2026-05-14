// Пакет archive реализует хранение данных согласно требованиям Закона Яровой:
// все сообщения, метаданные и информация о входах дублируются в изолированную БД
// с глубиной хранения 3 года.
package archive

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Service осуществляет дублирование данных в архивную БД.
type Service struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewService создаёт сервис архивирования.
func NewService(db *pgxpool.Pool, logger *zap.Logger) *Service {
	return &Service{db: db, logger: logger}
}

// ArchiveMessage сохраняет сообщение в архивную БД.
// Вызывается при каждом сохранении нового сообщения.
// Логи не содержат текст сообщения во избежание утечек.
func (s *Service) ArchiveMessage(ctx context.Context, messageID, senderID, chatID, body, mediaURL, senderIP string) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO archived_messages
		 (message_id, sender_id, chat_id, body, media_url, sender_ip, archived_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (message_id) DO NOTHING`,
		messageID, senderID, chatID, body, mediaURL, senderIP,
		time.Now(), time.Now().AddDate(3, 0, 0),
	)
	if err != nil {
		return fmt.Errorf("archive message: %w", err)
	}
	// Намеренно не логируем body — во избежание утечек через файлы журналов
	s.logger.Info("message archived",
		zap.String("messageID", messageID),
		zap.String("senderID", senderID))
	return nil
}

// ArchiveLogin сохраняет аудит-запись о входе пользователя.
func (s *Service) ArchiveLogin(ctx context.Context, userID, phone, ip string) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO audit_logins (user_id, phone, ip, logged_in_at, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		userID, phone, ip, time.Now(), time.Now().AddDate(3, 0, 0),
	)
	if err != nil {
		return fmt.Errorf("archive login: %w", err)
	}
	s.logger.Info("login archived",
		zap.String("userID", userID),
		zap.String("ip", ip))
	return nil
}
