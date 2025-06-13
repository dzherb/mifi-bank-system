package email_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dzherb/mifi-bank-system/pkg/email"
	"github.com/go-mail/mail/v2"
)

type MockDialer struct {
	CallsCount  int
	ReturnError error
}

func (d *MockDialer) DialAndSend(m ...*mail.Message) error {
	d.CallsCount++
	return d.ReturnError
}

func TestEmailSenderSuccess(t *testing.T) {
	dialerFactory := func() email.Dialer {
		return &MockDialer{
			ReturnError: nil,
		}
	}

	sender := email.NewSMTPSender(
		"smtp.test.com",
		587,
		"noreply@test.com",
		"123456",
	)
	sender.(*email.SMTPSender).SetDialerFactory(dialerFactory)

	err := sender.Send("user@test.com", "test", "test")
	if err != nil {
		t.Error(err)
	}
}

func TestEmailSenderFail(t *testing.T) {
	dialErr := fmt.Errorf("connection error")

	dialerFactory := func() email.Dialer {
		return &MockDialer{
			ReturnError: dialErr,
		}
	}

	sender := email.NewSMTPSender(
		"smtp.test.com",
		587,
		"noreply@test.com",
		"123456",
	)
	sender.(*email.SMTPSender).SetDialerFactory(dialerFactory)

	err := sender.Send("user@test.com", "test", "test")
	if err == nil {
		t.Errorf("expected error but got nil")
		return
	}

	if !errors.Is(err, dialErr) {
		t.Errorf("expected error %v but got %v", dialErr, err)
	}
}
