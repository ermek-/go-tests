package production_order

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
	Endpoint  = "/ProductionOrder/v1/ProductionOrders"
	Endpoints = Endpoint + "/"
)

type CreateRequest struct {
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

type Nomenclature struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	IsComposite *bool    `json:"is_composite"`
	ProdProcess int      `json:"prod_process"`
	Count       int      `json:"count"`
	Status      string   `json:"status"`
	Batches     []string `json:"batches"`
}
type CreateResponse struct {
	ID            int            `json:"id"`
	Number        *string        `json:"number"`
	DateGet       *string        `json:"date_get"`
	RequiredDate  *string        `json:"required_date"`
	DateComplete  *string        `json:"date_complite"`
	Priority      *int           `json:"priority"`
	Status        string         `json:"status"`
	Company       *int           `json:"company"`
	Customer      *string        `json:"customer"`
	ClientOrder   *int           `json:"client_order"`
	Nomenclatures []Nomenclature `json:"nomenclatures"`
}

func Create(t *testing.T, c *api.Client, endpoint string, body CreateRequest) (*http.Response, CreateResponse) {
	t.Helper()

	resp, err := c.Do(http.MethodPost, endpoint, body)
	require.NoError(t, err, "failed to create production order")

	b := helpers.ReadAllAndClose(t, resp)

	var po CreateResponse
	require.NoErrorf(t, json.Unmarshal(b, &po), "invalid JSON: %s", string(b))

	return resp, po
}

func Update(t *testing.T, c *api.Client, endpoint string, body CreateRequest, id int) *http.Response {
	t.Helper()

	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodPut, path, body)
	require.NoError(t, err, "failed to update production order")

	_ = resp.Body.Close()

	return resp
}

func Edit(t *testing.T, c *api.Client, endpoint string, body CreateRequest, id int) *http.Response {
	t.Helper()

	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodPatch, path, body)
	require.NoError(t, err, "failed to update production order")

	_ = resp.Body.Close()

	return resp
}

func GetList(t *testing.T, c *api.Client, endpoint string) (*http.Response, []byte, []CreateResponse) {
	t.Helper()

	path := fmt.Sprintf("%s/", endpoint)
	resp, err := c.Do(http.MethodGet, path, nil)
	require.NoError(t, err, "failed to GET production orders")

	body := helpers.ReadAllAndClose(t, resp)

	var list []CreateResponse
	require.NoErrorf(t, json.Unmarshal(body, &list), "invalid JSON: %s", string(body))

	return resp, body, list
}

func GetById(t *testing.T, c *api.Client, endpoint string, id int) (*http.Response, CreateResponse) {
	t.Helper()

	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodGet, path, nil)
	require.NoError(t, err, "failed to get production order")

	b := helpers.ReadAllAndClose(t, resp)

	var po CreateResponse
	require.NoErrorf(t, json.Unmarshal(b, &po), "invalid JSON: %s", string(b))

	return resp, po
}

func Delete(t *testing.T, c *api.Client, endpoint string, id int) {
	t.Helper()
	path := fmt.Sprintf("%s/%d/", endpoint, id)
	resp, err := c.Do(http.MethodDelete, path, nil)
	require.NoError(t, err, "Failed to delete production order")
	require.Equalf(t, http.StatusOK, resp.StatusCode, "Expected 200 OK, got %d", resp.StatusCode)
	_ = resp.Body.Close()
}
