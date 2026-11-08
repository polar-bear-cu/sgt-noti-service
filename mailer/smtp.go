package mailer

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net/smtp"
	"time"
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
	msg, err := buildMessage(m.from, to, subject, body)
	if err != nil {
		return err
	}
	return smtp.SendMail(m.addr, m.auth, m.from, []string{to}, msg)
}

func buildMessage(from, to, subject, body string) ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", subject))
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	w := quotedprintable.NewWriter(&b)
	if _, err := w.Write([]byte(body)); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
