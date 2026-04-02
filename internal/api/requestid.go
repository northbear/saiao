package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
)

const requestIDHeader = "X-Request-Id"

func requestIDFromRequest(r *http.Request) string {
	if r != nil {
		if requestID := normalizeRequestID(r.Header.Get(requestIDHeader)); requestID != "" {
			return requestID
		}
	}

	var buf [16]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return hex.EncodeToString(buf[:])
	}

	return "unknown"
}

func normalizeRequestID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > 128 {
		return ""
	}

	for _, r := range value {
		if r < 33 || r > 126 {
			return ""
		}
	}

	return value
}
