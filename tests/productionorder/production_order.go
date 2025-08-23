package productionorder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/internal/api"
	"go-api-tests/tests/helpers"

	"github.com/stretchr/testify/require"
)

const (
	Endpoint          = "/ProductionOrder/v1/ProductionOrders"
	EndpointWithSlash = Endpoint + "/"
)

type CreatePORequest struct {
	Number         *string `json:"number,omitempty"`
	Name           *string `json:"name,omitempty"`
	IsComposite    *bool   `json:"is_composite,omitempty"`
	ProdProcess    *int    `json:"prod_process,omitempty"`
	Count          *int    `json:"count,omitempty"`
	Status         *string `json:"status,omitempty"`
	TimePlanned    *string `json:"time_planned,omitempty"`
	AvailableCount *int    `json:"avaible_count,omitempty"`
	RemainedCount  *int    `json:"remained_count,omitempty"`
	ActualTime     *string `json:"actual_time,omitempty"`
	Material       *string `json:"material,omitempty"`
	Assortment     *string `json:"assortiment,omitempty"`
	Dimensions     *string `json:"dimensions,omitempty"`
	PassedCount    *int    `json:"passed_count,omitempty"`
	Dispatcher     *int    `json:"dispatcher,omitempty"`
}

type CreatePOResponse struct {
	ID            int           `json:"id"`
	Number        string        `json:"number"`
	DateGet       *string       `json:"date_get,omitempty"`
	RequiredDate  *string       `json:"required_date,omitempty"`
	DateComplete  *string       `json:"date_complite,omitempty"`
	Priority      *string       `json:"priority,omitempty"`
	Status        string        `json:"status"`
	Company       *string       `json:"company,omitempty"`
	Customer      *string       `json:"customer,omitempty"`
	ClientOrder   *string       `json:"client_order,omitempty"`
	Nomenclatures []interface{} `json:"nomenclatures,omitempty"`
}

func CreateProductionOrder(
	t *testing.T,
	c *api.Client,
	endpoint string,
	body CreatePORequest,
) (*http.Response, CreatePOResponse) {
	t.Helper()

	resp, err := c.Do(http.MethodPost, endpoint, body)
	require.NoError(t, err, "failed to create production order")

	b := helpers.ReadAllAndClose(t, resp)

	var po CreatePOResponse
	require.NoErrorf(t, json.Unmarshal(b, &po), "invalid JSON: %s", string(b))

	return resp, po
}

func GetProductionOrder(t *testing.T, c *api.Client, endpoint string, id int) (*http.Response, CreatePOResponse) {
	t.Helper()

	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodGet, path, nil)
	require.NoError(t, err, "failed to get production order")

	b := helpers.ReadAllAndClose(t, resp)

	var po CreatePOResponse
	require.NoErrorf(t, json.Unmarshal(b, &po), "invalid JSON: %s", string(b))

	return resp, po
}

func DeleteProductionOrder(t *testing.T, c *api.Client, endpoint string, id int) {
	t.Helper()
	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodDelete, path, nil)
	require.NoError(t, err, "Failed to delete production order")
	require.Equalf(t, http.StatusOK, resp.StatusCode, "Expected 200 OK, got %d", resp.StatusCode)
	_ = resp.Body.Close()
}
