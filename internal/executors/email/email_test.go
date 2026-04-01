package email

import (
	"bufio"
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func TestExecutorExecuteSendsRenderedEmail(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow opening a local listener")
		}
		t.Fatalf("listen smtp: %v", err)
	}
	defer listener.Close()

	messageCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go serveSMTP(listener, messageCh, errCh)

	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("split host port: %v", err)
	}

	executor := &Executor{}
	action := models.Action{
		To:              []string{"ops@example.com"},
		Cc:              []string{"lead@example.com"},
		SubjectTemplate: "Alert {{subject}}",
		BodyTemplate:    "Hello {{name}}",
	}
	identity := models.Identity{
		Email: "sender@example.com",
		Secrets: map[string]string{
			"smtp_host": host,
			"smtp_port": port,
		},
	}

	output, exitCode, err := executor.Execute(context.Background(), action, identity, map[string]any{
		"subject": "CPU",
		"name":    "team",
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != "email sent to 2 recipient(s)" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	select {
	case message := <-messageCh:
		if !strings.Contains(message, "Subject: Alert CPU") {
			t.Fatalf("missing subject in message: %q", message)
		}
		if !strings.Contains(message, "To: ops@example.com") {
			t.Fatalf("missing To header in message: %q", message)
		}
		if !strings.Contains(message, "Cc: lead@example.com") {
			t.Fatalf("missing Cc header in message: %q", message)
		}
		if !strings.Contains(message, "Hello team") {
			t.Fatalf("missing body in message: %q", message)
		}
	case err := <-errCh:
		t.Fatalf("smtp server error: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for smtp message")
	}
}

func TestExecutorExecuteReturnsTimeoutForCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	executor := &Executor{}
	_, exitCode, err := executor.Execute(ctx, models.Action{}, models.Identity{}, nil)
	if err == nil {
		t.Fatal("expected timeout")
	}
	if !errors.Is(err, saiaoerrors.ErrTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if exitCode != -1 {
		t.Fatalf("expected exit code -1, got %d", exitCode)
	}
}

func serveSMTP(listener net.Listener, messageCh chan<- string, errCh chan<- error) {
	conn, err := listener.Accept()
	if err != nil {
		errCh <- err
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writeSMTPLine(conn, "220 localhost ESMTP")

	inData := false
	var dataLines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		line = strings.TrimRight(line, "\r\n")

		if inData {
			if line == "." {
				writeSMTPLine(conn, "250 ok")
				messageCh <- strings.Join(dataLines, "\n")
				dataLines = nil
				inData = false
				continue
			}
			dataLines = append(dataLines, line)
			continue
		}

		switch {
		case strings.HasPrefix(line, "EHLO "), strings.HasPrefix(line, "HELO "):
			writeSMTPMultiline(conn, []string{"250-localhost", "250 ok"})
		case strings.HasPrefix(line, "MAIL FROM:"), strings.HasPrefix(line, "RCPT TO:"):
			writeSMTPLine(conn, "250 ok")
		case line == "DATA":
			writeSMTPLine(conn, "354 end with <CR><LF>.<CR><LF>")
			inData = true
		case line == "QUIT":
			writeSMTPLine(conn, "221 bye")
			return
		default:
			writeSMTPLine(conn, "250 ok")
		}
	}
}

func writeSMTPLine(conn net.Conn, line string) {
	_, _ = conn.Write([]byte(line + "\r\n"))
}

func writeSMTPMultiline(conn net.Conn, lines []string) {
	for _, line := range lines {
		writeSMTPLine(conn, line)
	}
}
