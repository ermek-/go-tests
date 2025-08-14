package helpers

import (
	"io"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func init() { gofakeit.Seed(0) }

func RandomNumber() int {
	return gofakeit.Number(1, 10000)
}

func ReadAllAndClose(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}
