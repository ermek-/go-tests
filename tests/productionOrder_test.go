package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// model
type ProductionOrder struct {
	ID            int     `json:"id"`
	Number        *string `json:"number"` // бывает null → указатель
	DateGet       *string `json:"date_get"`
	RequiredDate  *string `json:"required_date"`
	DateComplete  *string `json:"date_complite"` // оставил как в JSON
	Priority      *string `json:"priority"`
	Status        string  `json:"status"`
	Company       *string `json:"company"`
	Customer      *string `json:"customer"`
	ClientOrder   *string `json:"client_order"`
	Nomenclatures []any   `json:"nomenclatures"`
}

func TestCreateProductionOrder(t *testing.T) {
	t.Run("Create production order", func(t *testing.T) {
		// готовим запрос
		number := fmt.Sprintf("TEST-%d", RandomNumber())
		p := ProductionOrder{Number: &number}

		resp, err := client.Do(http.MethodPost, productionOrdersEndpoint, p)
		require.NoError(t, err, "Failed to create production order")
		require.Equalf(t, http.StatusCreated, resp.StatusCode, "Expected 201 Created, got %d", resp.StatusCode)

		// читаем тело ОДИН раз
		body := readAllAndClose(t, resp)

		// 1) проверяем наличие ключей на «сырых» данных
		var shape map[string]any
		require.NoErrorf(t, json.Unmarshal(body, &shape), "invalid JSON: %s", string(body))

		for _, k := range []string{
			"id", "number", "date_get", "required_date", "date_complite",
			"priority", "status", "company", "customer", "client_order",
			"nomenclatures",
		} {
			assert.Containsf(t, shape, k, "missing field %q in JSON body: %s", k, string(body))
		}

		// 2) разбираем в структуру
		var po ProductionOrder
		require.NoErrorf(t, json.Unmarshal(body, &po), "invalid JSON: %s", string(body))

		// осмысленные проверки
		assert.Equal(t, "Создан", po.Status)

		// 3) cleanup
		path := fmt.Sprintf("%s/%d/", productionOrderEndpoint, po.ID)

		deleteResp, err := client.Do(http.MethodDelete, path, nil)
		require.NoError(t, err, "Failed to delete production order")
		require.Equalf(t, http.StatusOK, deleteResp.StatusCode, "Expected 200 OK, got %d", deleteResp.StatusCode)
	})
}
