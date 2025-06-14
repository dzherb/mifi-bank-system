package notifications_test

import (
	"strings"
	"testing"

	"github.com/dzherb/mifi-bank-system/internal/services/notifications"
	"github.com/shopspring/decimal"
)

type MockEmailSender struct {
	To      string
	Subject string
	Body    string
}

func (s *MockEmailSender) Send(to, subject, body string) error {
	s.To = to
	s.Subject = subject
	s.Body = body

	return nil
}

func SetupMockEmailSender(t *testing.T) *MockEmailSender {
	ms := &MockEmailSender{}
	*notifications.EmailSender = ms

	t.Cleanup(func() {
		*notifications.EmailSender = nil
	})

	return ms
}

func TestSendPaymentSuccessEmail(t *testing.T) {
	mockSender := SetupMockEmailSender(t)

	to := "test@test.com"
	amount := decimal.NewFromFloat(150.5)

	err := notifications.SendPaymentSuccessEmail(to, amount)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if mockSender.To != to {
		t.Errorf("expected to %s, got %s", mockSender.To, to)
	}

	if mockSender.Subject != notifications.PaymentSuccessSubject {
		t.Errorf(
			"expected subject %s, got %s",
			mockSender.Subject,
			notifications.PaymentSuccessSubject,
		)
	}

	if !strings.Contains(mockSender.Body, amount.String()) {
		t.Errorf("expected body to contain %s", amount.String())
	}
}
