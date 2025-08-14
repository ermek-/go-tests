package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/internal/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ProductionOrder struct {
	ID            int     `json:"id"`
	Number        *string `json:"number"`
	DateGet       *string `json:"date_get"`
	RequiredDate  *string `json:"required_date"`
	DateComplete  *string `json:"date_complite"`
	Priority      *string `json:"priority"`
	Status        string  `json:"status"`
	Company       *string `json:"company"`
	Customer      *string `json:"customer"`
	ClientOrder   *string `json:"client_order"`
	Nomenclatures []any   `json:"nomenclatures"`
}

func CreateProductionOrder(t *testing.T, c *api.Client, endpoint string, body ProductionOrder) ProductionOrder {
	t.Helper()

	resp, err := c.Do(http.MethodPost, endpoint, body)
	require.NoError(t, err, "Failed to create production order")
	require.Equalf(t, http.StatusCreated, resp.StatusCode, "Expected 201 Created, got %d", resp.StatusCode)

	b := ReadAllAndClose(t, resp)

	var shape map[string]any
	require.NoErrorf(t, json.Unmarshal(b, &shape), "invalid JSON: %s", string(b))
	for _, k := range []string{
		"id", "number", "date_get", "required_date", "date_complite",
		"priority", "status", "company", "customer", "client_order",
		"nomenclatures",
	} {
		assert.Containsf(t, shape, k, "missing field %q in JSON body: %s", k, string(b))
	}

	var po ProductionOrder
	require.NoErrorf(t, json.Unmarshal(b, &po), "invalid JSON: %s", string(b))
	assert.Equal(t, "Создан", po.Status)
	return po
}

func DeleteProductionOrder(t *testing.T, c *api.Client, endpoint string, id int) {
	t.Helper()
	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodDelete, path, nil)
	require.NoError(t, err, "Failed to delete production order")
	require.Equalf(t, http.StatusOK, resp.StatusCode, "Expected 200 OK, got %d", resp.StatusCode)
	_ = resp.Body.Close()
}
