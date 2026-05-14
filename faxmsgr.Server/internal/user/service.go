// Пакет user управляет профилями пользователей.
package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Profile представляет профиль пользователя.
type Profile struct {
	ID          string `json:"id"`
	Phone       string `json:"phone"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

// UpdateProfileReq данные для обновления профиля.
type UpdateProfileReq struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// Service управляет профилями пользователей.
type Service struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// NewService создаёт сервис профилей.
func NewService(db *pgxpool.Pool, logger *zap.Logger) *Service {
	return &Service{db: db, logger: logger}
}

// GetProfile возвращает профиль пользователя по ID.
func (s *Service) GetProfile(ctx context.Context, userID string) (*Profile, error) {
	var p Profile
	err := s.db.QueryRow(ctx,
		`SELECT id, phone, COALESCE(display_name, ''), COALESCE(avatar_url, '')
		 FROM users WHERE id = $1`,
		userID,
	).Scan(&p.ID, &p.Phone, &p.DisplayName, &p.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

// UpdateProfile обновляет профиль пользователя.
func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileReq) (*Profile, error) {
	_, err := s.db.Exec(ctx,
		`UPDATE users SET display_name = $1, avatar_url = $2 WHERE id = $3`,
		req.DisplayName, req.AvatarURL, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return s.GetProfile(ctx, userID)
}
