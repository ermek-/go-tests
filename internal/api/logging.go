package api

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var redactedHeaders = map[string]struct{}{
	"authorization": {},
	"cookie":        {},
	"set-cookie":    {},
}

type logger interface{ Printf(string, ...any) }

type stdLogger struct {
	w  io.Writer
	mu sync.Mutex
}

func (l *stdLogger) Printf(format string, v ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.w, format, v...)
	if f, ok := l.w.(interface{ Flush() error }); ok {
		_ = f.Flush()
	}
}

func newWriterFromEnv() io.Writer {
	return bufio.NewWriter(os.Stdout)
}

type LoggingTransport struct {
	Base      http.RoundTripper
	LogBodies bool
	MaxBytes  int64
	L         logger
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	var reqBody []byte
	if t.LogBodies && req.Body != nil {
		limited := &io.LimitedReader{R: req.Body, N: t.MaxBytes}
		b, _ := io.ReadAll(limited)
		reqBody = b
		req.Body = io.NopCloser(io.MultiReader(bytes.NewReader(b), req.Body))
	}

	reqHeaders := headerString(req.Header)
	t.L.Printf("\n=== HTTP Request ===\n%s %s\n%s\n", req.Method, req.URL.String(), reqHeaders)
	if t.LogBodies && len(reqBody) > 0 {
		t.L.Printf("-- Request Body (<= %d bytes) --\n%s\n", t.MaxBytes, string(reqBody))
	}

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, err := base.RoundTrip(req)
	if err != nil {
		t.L.Printf("HTTP error: %v\n", err)
		return nil, err
	}
	dur := time.Since(start)

	var respBody []byte
	if t.LogBodies && resp.Body != nil {
		limited := &io.LimitedReader{R: resp.Body, N: t.MaxBytes}
		b, _ := io.ReadAll(limited)
		respBody = b
		resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(b), resp.Body))
	}

	respHeaders := headerString(resp.Header)
	t.L.Printf("\n=== HTTP Response ===\nStatus: %s in %s\n%s\n", resp.Status, dur.String(), respHeaders)
	if t.LogBodies && len(respBody) > 0 {
		t.L.Printf("-- Response Body (<= %d bytes) --\n%s\n", t.MaxBytes, string(respBody))
	}
	return resp, nil
}

func headerString(h http.Header) string {
	var b strings.Builder
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := strings.Join(h.Values(k), ", ")
		if _, ok := redactedHeaders[strings.ToLower(k)]; ok {
			v = "<redacted>"
		}
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString("\n")
	}
	return b.String()
}

func buildLoggingTransport() http.RoundTripper {
	if strings.TrimSpace(Env("LOG_HTTP", "")) == "" {
		return nil
	}
	logBodies := strings.TrimSpace(Env("LOG_HTTP_BODY", "")) != ""
	maxBytes := int64(8192)
	if v := strings.TrimSpace(Env("LOG_MAX_BYTES", "")); v != "" {
		if n, err := parseInt64(v); err == nil && n > 0 {
			maxBytes = n
		}
	}
	w := newWriterFromEnv()
	l := &stdLogger{w: w}
	return &LoggingTransport{LogBodies: logBodies, MaxBytes: maxBytes, L: l}
}

func parseInt64(s string) (int64, error) {
	var n int64
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int64(r-'0')
	}
	return n, nil
}
