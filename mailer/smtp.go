package mailer

import (
	"context"
	"fmt"
	"net/smtp"
)

// SMTPMailer sends plain-text mail over SMTP. Dev target is MailHog (no auth).
type SMTPMailer struct {
	addr string
	from string
	auth smtp.Auth
}

func NewSMTPMailer(host, port, from, username, password string) *SMTPMailer {
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}
	return &SMTPMailer{addr: host + ":" + port, from: from, auth: auth}
}

func (m *SMTPMailer) Send(_ context.Context, to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", m.from, to, subject, body)
	return smtp.SendMail(m.addr, m.auth, m.from, []string{to}, []byte(msg))
}
