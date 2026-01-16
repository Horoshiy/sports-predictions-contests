package channels

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailChannel struct {
	host     string
	port     string
	user     string
	password string
	from     string
	enabled  bool
}

func NewEmailChannel(host, port, user, password, from string, enabled bool) *EmailChannel {
	if !enabled || host == "" {
		return &EmailChannel{enabled: false}
	}
	return &EmailChannel{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
		enabled:  true,
	}
}

func (e *EmailChannel) Send(to, subject, body string) error {
	if !e.enabled {
		return nil
	}

	em := email.NewEmail()
	em.From = e.from
	em.To = []string{to}
	em.Subject = subject
	em.HTML = []byte(body)

	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	auth := smtp.PlainAuth("", e.user, e.password, e.host)

	if err := em.Send(addr, auth); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (e *EmailChannel) IsEnabled() bool {
	return e.enabled
}
