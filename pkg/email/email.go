package email

import (
	"crypto/tls"
	"fmt"

	"github.com/go-mail/mail/v2"
)

type Sender interface {
	Send(to, subject, body string) error
}

func NewSMTPSender(host string, port int, user, pass string) Sender {
	s := SMTPSender{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	}

	s.SetDialerFactory(s.createSMTPDialer)

	return &s
}

type SMTPSender struct {
	Host string
	Port int
	User string
	Pass string

	// We need this to mock Dialer in tests
	createDialer DialerFactory
}

type DialerFactory func() Dialer

type Dialer interface {
	DialAndSend(m ...*mail.Message) error
}

func (e *SMTPSender) Send(to, subject, body string) error {
	m := e.createMessage(to, subject, body)
	d := e.createDialer()

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("email sending failed: %w", err)
	}

	return nil
}

func (e *SMTPSender) SetDialerFactory(df DialerFactory) {
	e.createDialer = df
}

func (e *SMTPSender) createMessage(to, subject, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", e.User)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return m
}

func (e *SMTPSender) createSMTPDialer() Dialer { // coverage-ignore
	d := mail.NewDialer(e.Host, e.Port, e.User, e.Pass)
	d.TLSConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		ServerName:         e.Host,
		InsecureSkipVerify: false,
	}

	return d
}
