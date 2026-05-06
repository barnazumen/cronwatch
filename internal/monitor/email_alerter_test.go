package monitor

import (
	"net/smtp"
	"strings"
	"testing"
	"time"
)

func makeEmailAlerter(capturedMsg *[]byte, capturedTo *[]string, retErr error) *EmailAlerter {
	cfg := EmailAlerterConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		Username: "user@example.com",
		Password: "secret",
		From:     "cronwatch@example.com",
		To:       []string{"ops@example.com"},
	}
	a := NewEmailAlerter(cfg)
	a.send = func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		if capturedMsg != nil {
			*capturedMsg = msg
		}
		if capturedTo != nil {
			*capturedTo = to
		}
		return retErr
	}
	return a
}

func TestEmailAlerter_SendsOnFailure(t *testing.T) {
	var capturedMsg []byte
	var capturedTo []string

	a := makeEmailAlerter(&capturedMsg, &capturedTo, nil)

	event := AlertEvent{
		JobName:    "backup",
		Status:     "failed",
		Message:    "exit code 1",
		OccurredAt: time.Now(),
	}

	if err := a.Alert(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := string(capturedMsg)
	if !strings.Contains(body, "backup") {
		t.Errorf("expected job name in message, got: %s", body)
	}
	if !strings.Contains(body, "FAILED") {
		t.Errorf("expected status in subject, got: %s", body)
	}
	if !strings.Contains(body, "exit code 1") {
		t.Errorf("expected message body content, got: %s", body)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "ops@example.com" {
		t.Errorf("unexpected recipients: %v", capturedTo)
	}
}

func TestEmailAlerter_PropagatesError(t *testing.T) {
	expectedErr := fmt.Errorf("smtp connection refused")
	a := makeEmailAlerter(nil, nil, expectedErr)

	event := AlertEvent{
		JobName:    "sync",
		Status:     "missed",
		Message:    "deadline exceeded",
		OccurredAt: time.Now(),
	}

	if err := a.Alert(event); err == nil {
		t.Fatal("expected error but got nil")
	}
}
