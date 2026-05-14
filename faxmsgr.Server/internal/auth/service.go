// Пакет auth реализует аутентификацию через OTP и JWT.
package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	otpTTL          = 5 * time.Minute
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

// Service предоставляет методы аутентификации.
type Service struct {
	db        *pgxpool.Pool
	rdb       *redis.Client
	jwtSecret []byte
	logger    *zap.Logger
}

// NewService создаёт новый сервис аутентификации.
func NewService(db *pgxpool.Pool, rdb *redis.Client, jwtSecret string, logger *zap.Logger) *Service {
	return &Service{
		db:        db,
		rdb:       rdb,
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

// RequestCode генерирует OTP и сохраняет в Redis. В продакшене отправляет SMS.
func (s *Service) RequestCode(ctx context.Context, phone string) error {
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	key := otpKey(phone)
	if err := s.rdb.Set(ctx, key, code, otpTTL).Err(); err != nil {
		return fmt.Errorf("store otp: %w", err)
	}

	// TODO: отправить SMS через шлюз (СМSC, МТС и т.д.)
	s.logger.Info("OTP generated (SMS sending not implemented)",
		zap.String("phone", phone))

	return nil
}

// VerifyCode проверяет OTP и выдаёт пару токенов.
func (s *Service) VerifyCode(ctx context.Context, phone, code string) (accessToken, refreshToken string, err error) {
	key := otpKey(phone)
	stored, err := s.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", "", errors.New("OTP not found or expired")
	}
	if err != nil {
		return "", "", fmt.Errorf("redis get otp: %w", err)
	}
	if stored != code {
		return "", "", errors.New("invalid OTP")
	}
	s.rdb.Del(ctx, key) //nolint:errcheck

	// Найти или создать пользователя
	userID, err := s.findOrCreateUser(ctx, phone)
	if err != nil {
		return "", "", fmt.Errorf("find or create user: %w", err)
	}

	// Выпустить Access JWT
	accessToken, err = s.issueAccessToken(userID)
	if err != nil {
		return "", "", fmt.Errorf("issue access token: %w", err)
	}

	// Выпустить Refresh JWT и сохранить хеш в Redis
	refreshToken, err = s.issueRefreshToken(ctx, userID)
	if err != nil {
		return "", "", fmt.Errorf("issue refresh token: %w", err)
	}

	// Аудит-лог: фиксируем вход (без текста сообщений)
	s.logger.Info("user authenticated",
		zap.String("userID", userID),
		zap.String("phone", phone))

	return accessToken, refreshToken, nil
}

// Refresh обновляет пару токенов по refresh-токену.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (newAccess, newRefresh string, err error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("parse refresh token: %w", err)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	// Проверить, что refresh-токен есть в Redis (не отозван)
	rkey := refreshKey(userID)
	storedHash, err := s.rdb.Get(ctx, rkey).Result()
	if errors.Is(err, redis.Nil) {
		return "", "", errors.New("refresh token revoked")
	}
	if err != nil {
		return "", "", fmt.Errorf("redis get refresh: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refreshToken)); err != nil {
		return "", "", errors.New("refresh token mismatch")
	}

	newAccess, err = s.issueAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	newRefresh, err = s.issueRefreshToken(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

// --- внутренние методы ---

func (s *Service) findOrCreateUser(ctx context.Context, phone string) (string, error) {
	var id string
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (phone) VALUES ($1)
		 ON CONFLICT (phone) DO UPDATE SET phone = EXCLUDED.phone
		 RETURNING id`,
		phone,
	).Scan(&id)
	return id, err
}

func (s *Service) issueAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(accessTokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) issueRefreshToken(ctx context.Context, userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"exp":  time.Now().Add(refreshTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	// Хешируем refresh-токен перед сохранением в Redis
	hash, err := bcrypt.GenerateFromPassword([]byte(tokenStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	rkey := refreshKey(userID)
	if err := s.rdb.Set(ctx, rkey, string(hash), refreshTokenTTL).Err(); err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (s *Service) parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// generateOTP генерирует криптографически случайный числовой код заданной длины.
func generateOTP(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", digits, n), nil
}

func otpKey(phone string) string     { return "otp:" + phone }
func refreshKey(userID string) string { return "refresh:" + userID }
