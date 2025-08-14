package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-api-tests/tests/helpers"
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
		number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
		order := ProductionOrder{Number: &number}

		resp, err := client.Do(http.MethodPost, productionOrdersEndpoint, order)
		require.NoError(t, err, "Failed to create production order")
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var created ProductionOrder
		body := helpers.ReadAllAndClose(t, resp)
		require.NoError(t, json.Unmarshal(body, &created))

		deleteResp, err := client.Do(
			http.MethodDelete,
			fmt.Sprintf("%s/%d/", productionOrderEndpoint, created.ID),
			nil,
		)
		require.NoError(t, err, "Failed to delete production order")
		require.Equal(t, http.StatusOK, deleteResp.StatusCode)
		_ = deleteResp.Body.Close()
	})
}
