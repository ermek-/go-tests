package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
)

type Product struct {
	ID    any    `json:"id,omitempty"`
	Code  string `json:"code"`
	Price int    `json:"price"`
}

func readJSON(resp *http.Response, v any, t *testing.T) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, string(b))
	}
}

func TestCRUD_Product(t *testing.T) {
	var createdID any

	t.Run("Create", func(t *testing.T) {
		code := fmt.Sprintf("code-%d", rand.Intn(1000000))
		price := rand.Intn(1000) + 1 // от 1 до 1000

		p := Product{Code: code, Price: price}

		resp, err := client.Do(http.MethodPost, productEndpoint, p)
		if err != nil {
			t.Fatalf("POST: %v", err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var m map[string]any
		readJSON(resp, &m, t)
		id, ok := m["id"]
		if !ok {
			t.Fatalf("no id in response: %v", m)
		}
		createdID = id
	})

	t.Run("GetByID", func(t *testing.T) {
		path := fmt.Sprintf("%s/%v", productsEndpoint, createdID)
		resp, err := client.Do(http.MethodGet, path, nil)
		if err != nil {
			t.Fatalf("GET by id: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var p Product
		readJSON(resp, &p, t)
		if p.Code == "" {
			t.Fatalf("empty code in product: %+v", p)
		}
	})

	t.Run("Update", func(t *testing.T) {
		path := fmt.Sprintf("%s/%v", productsEndpoint, createdID)
		upd := map[string]any{"price": 150}
		resp, err := client.Do(http.MethodPut, path, upd)
		if err != nil {
			t.Fatalf("PUT: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
	})

	t.Run("List", func(t *testing.T) {
		resp, err := client.Do(http.MethodGet, productsEndpoint, nil)
		if err != nil {
			t.Fatalf("GET list: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		var list []map[string]any
		readJSON(resp, &list, t)
		if len(list) == 0 {
			t.Fatalf("expected non-empty list")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		path := fmt.Sprintf("%s/%v", productsEndpoint, createdID)
		resp, err := client.Do(http.MethodDelete, path, nil)
		if err != nil {
			t.Fatalf("DELETE: %v", err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
	})
}
