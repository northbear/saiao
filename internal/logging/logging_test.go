package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewWithWriterJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithWriter(&buf, "json", "info")

	logger.Info("hello", "component", "test")

	output := buf.String()
	if !strings.Contains(output, `"msg":"hello"`) {
		t.Fatalf("expected JSON log message, got %q", output)
	}
	if !strings.Contains(output, `"component":"test"`) {
		t.Fatalf("expected JSON field, got %q", output)
	}
}

func TestNewWithWriterTextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithWriter(&buf, "text", "info")

	logger.Info("hello", "component", "test")

	output := buf.String()
	if !strings.Contains(output, "msg=hello") {
		t.Fatalf("expected text log message, got %q", output)
	}
	if !strings.Contains(output, "component=test") {
		t.Fatalf("expected text field, got %q", output)
	}
}

func TestNewWithWriterHonorsLevel(t *testing.T) {
	var infoBuf bytes.Buffer
	infoLogger := NewWithWriter(&infoBuf, "json", "info")
	infoLogger.Debug("debug-hidden")
	if infoBuf.Len() != 0 {
		t.Fatalf("expected debug log to be suppressed at info level, got %q", infoBuf.String())
	}

	var debugBuf bytes.Buffer
	debugLogger := NewWithWriter(&debugBuf, "json", "debug")
	debugLogger.Debug("debug-visible")
	if !strings.Contains(debugBuf.String(), `"msg":"debug-visible"`) {
		t.Fatalf("expected debug log to be emitted, got %q", debugBuf.String())
	}
}

func TestDiscardLoggerAcceptsWrites(t *testing.T) {
	logger := Discard()
	logger.Info("hello", "component", "test")
	logger.Debug("debug")
}
