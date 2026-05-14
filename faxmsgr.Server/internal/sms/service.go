// Пакет sms предоставляет интерфейс для отправки SMS-сообщений.
// Реальная интеграция с оператором (СМSC, МТС и т.д.) реализуется в конкретной реализации.
package sms

import (
	"context"
	"fmt"
)

// Sender интерфейс SMS-отправителя.
type Sender interface {
	Send(ctx context.Context, phone, text string) error
}

// StubSender заглушка SMS-отправителя для разработки.
// В продакшен-среде заменяется на реальный шлюз.
type StubSender struct{}

// Send выводит SMS в stdout вместо реальной отправки.
func (s *StubSender) Send(_ context.Context, phone, text string) error {
	fmt.Printf("[SMS STUB] To: %s | Text: %s\n", phone, text)
	return nil
}
