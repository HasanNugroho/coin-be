package utils

import (
	"context"
	"fmt"
	"net"
	"net/mail"
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
	fmt.Printf("[SMTP] Starting Send\n  HOST: %s\n  PORT: %s\n  USERNAME: %s\n  FROM: %s\n  TO: %s\n  PASSWORD set: %v\n",
		m.host, m.port, m.username, m.from, to, m.password != "")

	if m.host == "" {
		fmt.Printf("[SMTP] No host set, printing email:\nSubject: %s\nBody: %s\n", subject, body)
		return nil
	}

	if to == "" {
		return fmt.Errorf("SMTP Send called with empty recipient address")
	}

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	fmt.Printf("[SMTP] Dialing %s...\n", addr)

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Printf("[SMTP] Connection failed: %v\n", err)
		return fmt.Errorf("SMTP connection failed: %w", err)
	}
	fmt.Printf("[SMTP] TCP connection established\n")

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		conn.Close()
		fmt.Printf("[SMTP] Failed to create client: %v\n", err)
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()
	fmt.Printf("[SMTP] SMTP client created\n")

	if m.username != "" && m.password != "" {
		fmt.Printf("[SMTP] Authenticating as %s...\n", m.username)
		auth := plainAuth{username: m.username, password: m.password}
		if err := client.Auth(auth); err != nil {
			fmt.Printf("[SMTP] Auth failed: %v\n", err)
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
		fmt.Printf("[SMTP] Auth successful\n")
	} else {
		fmt.Printf("[SMTP] Skipping auth (no credentials)\n")
	}

	plainFrom := m.from
	parsed, err := mail.ParseAddress(m.from)
	if err == nil {
		plainFrom = parsed.Address
	}
	fmt.Printf("[SMTP] MAIL FROM: %s (parsed from: %s)\n", plainFrom, m.from)
	if err := client.Mail(plainFrom); err != nil {
		fmt.Printf("[SMTP] MAIL FROM failed: %v\n", err)
		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
	}
	fmt.Printf("[SMTP] MAIL FROM accepted\n")

	plainTo := to
	parsedTo, err := mail.ParseAddress(to)
	if err == nil {
		plainTo = parsedTo.Address
	}
	fmt.Printf("[SMTP] RCPT TO: %s (parsed from: %s)\n", plainTo, to)
	if err := client.Rcpt(plainTo); err != nil {
		fmt.Printf("[SMTP] RCPT TO failed: %v\n", err)
		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
	}
	fmt.Printf("[SMTP] RCPT TO accepted\n")

	fmt.Printf("[SMTP] Opening DATA writer...\n")
	w, err := client.Data()
	if err != nil {
		fmt.Printf("[SMTP] DATA failed: %v\n", err)
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.from, to, subject, body,
	)
	fmt.Printf("[SMTP] Writing message (%d bytes)...\n", len(msg))
	if _, err := fmt.Fprint(w, msg); err != nil {
		fmt.Printf("[SMTP] Write failed: %v\n", err)
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	// Close writer secara eksplisit untuk trigger pengiriman ke server
	if err := w.Close(); err != nil {
		fmt.Printf("[SMTP] Close (flush) failed: %v\n", err)
		return fmt.Errorf("SMTP flush failed: %w", err)
	}
	fmt.Printf("[SMTP] Message flushed to server\n")

	// QUIT dan tangkap response akhir dari server
	if err := client.Quit(); err != nil {
		fmt.Printf("[SMTP] QUIT response: %v\n", err)
	} else {
		fmt.Printf("[SMTP] QUIT accepted, server selesai memproses\n")
	}

	fmt.Printf("[SMTP] Email sent successfully to %s\n", to)
	return nil
}

// func (m *smtpMailer) Send(ctx context.Context, to, subject, body string) error {
// 	if m.host == "" {
// 		fmt.Printf("DEBUG: Sending email to %s\nSubject: %s\nBody: %s\n", to, subject, body)
// 		return nil
// 	}

// 	addr := fmt.Sprintf("%s:%s", m.host, m.port)

// 	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
// 	if err != nil {
// 		return fmt.Errorf("SMTP connection failed: %w", err)
// 	}

// 	client, err := smtp.NewClient(conn, m.host)
// 	if err != nil {
// 		conn.Close()
// 		return fmt.Errorf("failed to create SMTP client: %w", err)
// 	}
// 	defer client.Close()

// 	if m.username != "" && m.password != "" {
// 		auth := plainAuth{username: m.username, password: m.password}
// 		if err := client.Auth(auth); err != nil {
// 			return fmt.Errorf("SMTP auth failed: %w", err)
// 		}
// 	}

// 	if err := client.Mail(m.from); err != nil {
// 		return fmt.Errorf("SMTP MAIL FROM failed: %w", err)
// 	}
// 	if err := client.Rcpt(to); err != nil {
// 		return fmt.Errorf("SMTP RCPT TO failed: %w", err)
// 	}

// 	w, err := client.Data()
// 	if err != nil {
// 		return fmt.Errorf("SMTP DATA failed: %w", err)
// 	}
// 	defer w.Close()

// 	msg := fmt.Sprintf(
// 		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
// 		m.from, to, subject, body,
// 	)
// 	if _, err := fmt.Fprint(w, msg); err != nil {
// 		return fmt.Errorf("SMTP write failed: %w", err)
// 	}

// 	return nil
// }
