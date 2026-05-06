package monitor

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// EmailAlerterConfig holds SMTP configuration for sending alert emails.
type EmailAlerterConfig struct {
	SMTPHost   string
	SMTPPort   int
	Username   string
	Password   string
	From       string
	To         []string
}

// EmailAlerter sends job alerts via email.
type EmailAlerter struct {
	cfg  EmailAlerterConfig
	send func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// NewEmailAlerter creates an EmailAlerter with the given SMTP configuration.
func NewEmailAlerter(cfg EmailAlerterConfig) *EmailAlerter {
	return &EmailAlerter{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

// Alert sends an email notification for the given AlertEvent.
func (e *EmailAlerter) Alert(event AlertEvent) error {
	subject := fmt.Sprintf("[cronwatch] %s: job %q", strings.ToUpper(event.Status), event.JobName)
	body := fmt.Sprintf(
		"Job:     %s\nStatus:  %s\nTime:    %s\nMessage: %s\n",
		event.JobName,
		event.Status,
		event.OccurredAt.Format(time.RFC1123),
		event.Message,
	)

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n",
		e.cfg.From,
		strings.Join(e.cfg.To, ", "),
		subject,
	)

	msg := []byte(header + body)
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)

	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.SMTPHost)
	}

	return e.send(addr, auth, e.cfg.From, e.cfg.To, msg)
}
