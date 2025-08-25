package production_order

import (
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	c := helpers.TestClient(t)
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())

	req := CreateRequest{Number: &number}
	resp, po := Create(t, c, Endpoints, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")
	require.Equal(t, number, po.Number, "expected production order number to match")
	require.NotZero(t, po.ID, "server must return non-zero id")
	require.Equal(t, "Создан", po.Status, "expected production order status to match")

	getResp, got := GetById(t, c, Endpoint, po.ID)
	require.Equal(t, http.StatusOK, getResp.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, po.ID, got.ID, "GET should return the same ID")
	assert.Equal(t, number, got.Number, "GET should return the same Number")
	assert.Equal(t, "Создан", got.Status, "status should match after create")

	t.Cleanup(func() { Delete(t, c, Endpoint, po.ID) })
}

func TestGetList(t *testing.T) {
	c := helpers.TestClient(t)

	resp, body, list := GetList(t, c, Endpoint)
	require.Equal(t, http.StatusOK, resp.StatusCode, "expected 200 OK")
	require.NotEmpty(t, list, "list must not be empty")
	helpers.AssertAllObjectsHaveKeysJSON(t, body,
		"id", "number", "date_get", "required_date", "date_complite",
		"priority", "status", "company", "customer", "client_order", "nomenclatures",
	)
}

func TestUpdate(t *testing.T) {
	c := helpers.TestClient(t)
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())

	req := CreateRequest{Number: &number}
	resp, po := Create(t, c, Endpoints, req)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")

	updNumber := fmt.Sprintf("%d", helpers.RandomNumber())
	updReq := CreateRequest{Number: &updNumber}
	updResp := Update(t, c, Endpoint, updReq, po.ID)
	require.Equal(t, http.StatusOK, updResp.StatusCode, "expected 200 OK on update")

	getResp, got := GetById(t, c, Endpoint, po.ID)
	require.Equal(t, http.StatusOK, getResp.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, updNumber, *got.Number, "GET should return the same Number")

	editNumber := fmt.Sprintf("%d", helpers.RandomNumber())
	editReq := CreateRequest{Number: &editNumber}
	editResp := Edit(t, c, Endpoint, editReq, po.ID)
	require.Equal(t, http.StatusOK, editResp.StatusCode, "expected 200 OK on edit")

	getResp2, got := GetById(t, c, Endpoint, po.ID)
	require.Equal(t, http.StatusOK, getResp2.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, editNumber, *got.Number, "GET should return the same Number")

	t.Cleanup(func() { Delete(t, c, Endpoint, po.ID) })
}
