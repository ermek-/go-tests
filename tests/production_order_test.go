package tests

import (
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/tests/helpers"
	"go-api-tests/tests/productionorder"

	"github.com/stretchr/testify/assert"
)

func TestCreateProductionOrder(t *testing.T) {
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
	req := productionorder.CreatePORequest{Number: &number}

	resp, po := productionorder.CreateProductionOrder(t, client, productionorder.EndpointWithSlash, req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")
	assert.Equal(t, number, po.Number, "expected production order number to match")
	assert.NotZero(t, po.ID, "server must return non-zero id")
	assert.Equal(t, "Создан", po.Status, "expected production order status to match")

	getResp, got := productionorder.GetProductionOrder(t, client, productionorder.Endpoint, po.ID)
	assert.Equal(t, http.StatusOK, getResp.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, po.ID, got.ID, "GET should return the same ID")
	assert.Equal(t, number, got.Number, "GET should return the same Number")
	assert.Equal(t, "Создан", got.Status, "status should match after create")

	productionorder.DeleteProductionOrder(t, client, productionorder.Endpoint, po.ID)
}
