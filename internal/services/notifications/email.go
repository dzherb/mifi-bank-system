package notifications

import (
	"fmt"

	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/dzherb/mifi-bank-system/pkg/email"
	"github.com/shopspring/decimal"
)

func Init(cfg *config.Config) { // coverage-ignore
	emailSender = email.NewSMTPSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)
}

var emailSender email.Sender

func activeSender() (email.Sender, error) {
	if emailSender == nil {
		return nil, fmt.Errorf("email sender not initialized")
	}

	return emailSender, nil
}

const paymentSuccessSubject = "Платеж успешно проведен"

func SendPaymentSuccessEmail(userEmail string, amount decimal.Decimal) error {
	sender, err := activeSender()
	if err != nil {
		return err
	}

	content := fmt.Sprintf(`
        <h1>Спасибо за оплату!</h1>
        <p>Сумма: <strong>%s RUB</strong></p>
        <small>Это автоматическое уведомление</small>`,
		amount.String(),
	)

	err = sender.Send(userEmail, paymentSuccessSubject, content)
	if err != nil {
		return err
	}

	return nil
}
