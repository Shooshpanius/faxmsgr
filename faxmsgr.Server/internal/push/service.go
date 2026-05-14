// Пакет push предоставляет интерфейс для отправки push-уведомлений (FCM и аналоги).
package push

import (
	"context"
	"fmt"
)

// Notifier интерфейс сервиса push-уведомлений.
type Notifier interface {
	Notify(ctx context.Context, deviceToken, title, body string) error
}

// StubNotifier заглушка push-уведомлений для разработки.
type StubNotifier struct{}

// Notify выводит уведомление в stdout вместо реальной отправки.
func (n *StubNotifier) Notify(_ context.Context, deviceToken, title, body string) error {
	fmt.Printf("[PUSH STUB] Token: %s | %s: %s\n", deviceToken, title, body)
	return nil
}
