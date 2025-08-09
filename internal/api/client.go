package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	BaseURL   string
	Token     string
	HTTP      *http.Client
	UserAgent string
}

func NewClient(baseURL string) *Client {
	c := &Client{
		BaseURL:   strings.TrimRight(baseURL, "/"),
		UserAgent: "go-api-tests/1.0",
		HTTP:      &http.Client{Timeout: 15 * time.Second},
	}
	if tr := buildLoggingTransport(); tr != nil {
		c.HTTP.Transport = tr // set only when enabled to avoid typed-nil in interface
	}
	return c
}

func (c *Client) Authenticate(authEndpoint, username, password string) error {
	endpoint := strings.TrimSpace(authEndpoint)
	if endpoint == "" {
		return nil // auth not required
	}
	url := c.BaseURL + endpoint

	payload := map[string]string{
		"username": username,
		"password": password,
	}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// 🔹 Читаем значение из ENV и ставим заголовок, если не пусто
	if apiKey := Env("X-CSRFTOKEN", ""); apiKey != "" {
		req.Header.Set("X-CSRFTOKEN", apiKey)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return fmt.Errorf("decode auth response: %w", err)
	}
	for _, k := range []string{"access", "access_token", "jwt"} {
		if v, ok := m[k]; ok {
			if s, ok2 := v.(string); ok2 && s != "" {
				c.Token = s
				return nil
			}
		}
	}
	return fmt.Errorf("no token found in response")
}

func (c *Client) Do(method, path string, body any) (*http.Response, error) {
	url := c.BaseURL + path
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "JWT "+c.Token)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	return c.HTTP.Do(req)
}

func Env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}
