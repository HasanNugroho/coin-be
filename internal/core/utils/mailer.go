package utils

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"
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

type plainAuth struct {
	username, password string
}

func (a plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	credentials := fmt.Sprintf("\x00%s\x00%s", a.username, a.password)
	return "PLAIN", []byte(credentials), nil
}

func (a plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("unexpected server challenge")
	}
	return nil, nil
}

func (m *smtpMailer) Send(ctx context.Context, to, subject, body string) error {
	if m.host == "" {
		fmt.Printf("DEBUG: Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("SMTP connection failed: %w", err)
	}

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if m.username != "" && m.password != "" {
		auth := plainAuth{username: m.username, password: m.password}
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err := client.Mail(m.from); err != nil {
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}
	defer w.Close()

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)
	if _, err := fmt.Fprint(w, msg); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	return nil
}
