package email

import (
	"crypto/tls"
	"fmt"

	"github.com/go-mail/mail/v2"
)

type Sender struct {
	Host string
	Port int
	User string
	Pass string
}

func (e *Sender) Send(to, subject, body string) error {
	m := e.createMessage(to, subject, body)
	d := e.createDialer()

	if err := e.send(d, m); err != nil {
		return err
	}

	return nil
}

func (e *Sender) createMessage(to, subject, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", e.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return m
}

func (e *Sender) createDialer() *mail.Dialer {
	d := mail.NewDialer(e.Host, e.Port, e.User, e.Pass)
	d.TLSConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		ServerName:         e.Host,
		InsecureSkipVerify: false,
	}

	return d
}

func (e *Sender) send(d *mail.Dialer, m *mail.Message) error {
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("email sending failed: %w", err)
	}

	return nil
}
