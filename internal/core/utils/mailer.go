package utils

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

type smtpMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewMailer() Mailer {
	return &smtpMailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (m *smtpMailer) Send(ctx context.Context, to, subject, body string) error {
	// If no SMTP config, log and return (for development)
	if m.host == "" {
		fmt.Printf("DEBUG: Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	return smtp.SendMail(addr, auth, m.from, []string{to}, msg)
}
