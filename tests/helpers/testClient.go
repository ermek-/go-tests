package helpers

import (
	"sync"
	"testing"

	"go-api-tests/internal/api"
	"go-api-tests/internal/env"
)

var (
	once    sync.Once
	client  *api.Client
	initErr error
)

func TestClient(t *testing.T) *api.Client {
	t.Helper()

	once.Do(func() {
		_ = env.LoadDotEnv(".env")

		baseURL := api.Env("BASE_URL", "http://api.dev.pkmt.tech")
		authEndpoint := api.Env("AUTH_ENDPOINT", "")
		username := api.Env("USERNAME", "")
		password := api.Env("PASSWORD", "")

		c := api.NewClient(baseURL)

		if authEndpoint != "" {
			if err := c.Authenticate(authEndpoint, username, password); err != nil {
				initErr = err
				return
			}
		}
		client = c
	})

	if initErr != nil {
		t.Fatalf("failed to init test client: %v", initErr)
	}
	if client == nil {
		t.Fatalf("test client is nil")
	}
	return client
}
