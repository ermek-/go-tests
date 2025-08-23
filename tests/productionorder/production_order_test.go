package productionorder

import (
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/tests/helpers"

	"github.com/stretchr/testify/assert"
)

func TestCreateProductionOrder(t *testing.T) {
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
	req := CreatePORequest{Number: &number}

	resp, po := CreateProductionOrder(t, client, EndpointWithSlash, req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")
	assert.Equal(t, number, po.Number, "expected production order number to match")
	assert.NotZero(t, po.ID, "server must return non-zero id")
	assert.Equal(t, "Создан", po.Status, "expected production order status to match")

	getResp, got := GetProductionOrder(t, client, Endpoint, po.ID)
	assert.Equal(t, http.StatusOK, getResp.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, po.ID, got.ID, "GET should return the same ID")
	assert.Equal(t, number, got.Number, "GET should return the same Number")
	assert.Equal(t, "Создан", got.Status, "status should match after create")

	DeleteProductionOrder(t, client, Endpoint, po.ID)
}
