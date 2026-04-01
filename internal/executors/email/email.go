package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
	"saiao/internal/render"
)

type Executor struct{}

func (e *Executor) Execute(ctx context.Context, action models.Action, identity models.Identity, input map[string]any) (string, int, error) {
	if err := ctx.Err(); err != nil {
		return "", -1, fmt.Errorf("%w: %v", saiaoerrors.ErrTimeout, err)
	}

	subject, err := render.Template(action.SubjectTemplate, identity, input)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
	}

	body, err := render.Template(action.BodyTemplate, identity, input)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
	}

	recipients, err := renderRecipients(identity, input, append(append([]string{}, action.To...), append(action.Cc, action.Bcc...)...))
	if err != nil {
		return "", 0, err
	}

	host := strings.TrimSpace(identity.Secrets["smtp_host"])
	port := strings.TrimSpace(identity.Secrets["smtp_port"])
	username := identity.Secrets["smtp_username"]
	password := identity.Secrets["smtp_password"]
	from := identity.Secrets["from_email"]
	if from == "" {
		from = identity.Email
	}

	if host == "" || port == "" || from == "" {
		return "", 0, fmt.Errorf("%w: email identity is missing smtp_host, smtp_port, or from_email", saiaoerrors.ErrExecutionFailed)
	}

	message := buildMessage(from, action, recipients, subject, body)
	address := host + ":" + port

	var auth smtp.Auth
	if username != "" || password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	if err := smtp.SendMail(address, auth, from, recipients, []byte(message)); err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
	}

	return fmt.Sprintf("email sent to %d recipient(s)", len(recipients)), 0, nil
}

func renderRecipients(identity models.Identity, input map[string]any, recipients []string) ([]string, error) {
	rendered := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		value, err := render.Template(recipient, identity, input)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
		}
		value = strings.TrimSpace(value)
		if value != "" {
			rendered = append(rendered, value)
		}
	}
	return rendered, nil
}

func buildMessage(from string, action models.Action, recipients []string, subject, body string) string {
	headers := []string{
		"From: " + from,
		"To: " + strings.Join(recipients[:len(action.To)], ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}

	if len(action.Cc) > 0 {
		ccStart := len(action.To)
		ccEnd := ccStart + len(action.Cc)
		headers = append(headers[:2], append([]string{"Cc: " + strings.Join(recipients[ccStart:ccEnd], ", ")}, headers[2:]...)...)
	}

	return strings.Join(headers, "\r\n")
}
