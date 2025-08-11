package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Product struct {
	ID    any    `json:"id,omitempty"`
	Code  string `json:"code"`
	Price int    `json:"price"`
}

func TestCRUD_Product(t *testing.T) {
	var createdID any

	t.Run("Create", func(t *testing.T) {
		code := fmt.Sprintf("code-%d", RandomNumber())
		price := RandomNumber()
		p := Product{Code: code, Price: price}

		resp, err := client.Do(http.MethodPost, nomenclaturesEndpoint, p)
		require.NoError(t, err, "POST request failed")

		require.Truef(t,
			resp.StatusCode == http.StatusCreated,
			"unexpected status: %d", resp.StatusCode,
		)

		body := readAllAndClose(t, resp)

		var m map[string]any
		require.NoErrorf(t, json.Unmarshal(body, &m), "invalid JSON: %s", string(body))

		for _, k := range []string{
			"id", "name", "is_composite", "nomenclature_in_process_id", "prod_process",
			"count", "status", "time_planned", "avaible_count", "remained_count",
			"actual_time", "material", "assortiment", "dimensions", "passed_count", "dispatcher",
		} {
			assert.Containsf(t, m, k, "missing field %q in JSON body: %s", k, string(body))
		}

		id, ok := m["id"]
		require.Truef(t, ok, "no id in response: %v", m)
		createdID = id
	})

	t.Run("GetByID", func(t *testing.T) {
		require.NotNil(t, createdID, "createdID should be set by Create")

		path := fmt.Sprintf("%s/%v", nomenclatureEndpoint, createdID)
		resp, err := client.Do(http.MethodGet, path, nil)
		require.NoError(t, err, "GET by id failed")
		require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status")

		var p Product
		// читаем тело 1 раз, затем decode
		body := readAllAndClose(t, resp)
		require.NoErrorf(t, json.Unmarshal(body, &p), "invalid JSON: %s", string(body))
		assert.NotEmpty(t, p.Code, "product code should not be empty")
	})

	t.Run("Update", func(t *testing.T) {
		require.NotNil(t, createdID)

		path := fmt.Sprintf("%s/%v", nomenclatureEndpoint, createdID)
		upd := map[string]any{"price": 150}
		resp, err := client.Do(http.MethodPut, path, upd)
		require.NoError(t, err, "PUT failed")
		require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status")
		_ = resp.Body.Close()
	})

	t.Run("List", func(t *testing.T) {
		resp, err := client.Do(http.MethodGet, nomenclatureEndpoint, nil)
		require.NoError(t, err, "GET list failed")
		require.Equal(t, http.StatusOK, resp.StatusCode, "unexpected status")

		var list []map[string]any
		body := readAllAndClose(t, resp)
		require.NoErrorf(t, json.Unmarshal(body, &list), "invalid list JSON: %s", string(body))
		require.Greater(t, len(list), 0, "expected non-empty list")
	})

	t.Run("Delete", func(t *testing.T) {
		require.NotNil(t, createdID)

		path := fmt.Sprintf("%s/%v/", nomenclatureEndpoint, createdID) // если слэш лишний — убери
		resp, err := client.Do(http.MethodDelete, path, nil)
		require.NoError(t, err, "DELETE failed")
		require.Truef(t,
			resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent,
			"unexpected status: %d", resp.StatusCode,
		)
		_ = resp.Body.Close()
	})
}
