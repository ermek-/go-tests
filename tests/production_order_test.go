package tests

import (
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/tests/helpers"

	"github.com/stretchr/testify/assert"
)

func TestCreateProductionOrder(t *testing.T) {
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
	req := helpers.CreatePORequest{Number: &number}

	resp, po := helpers.CreateProductionOrder(t, client, productionOrdersEndpoint, req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")
	assert.Equal(t, number, po.Number)
	assert.NotZero(t, po.ID, "server must return non-zero id")
	assert.Equal(t, "Создан", po.Status)

	helpers.DeleteProductionOrder(t, client, productionOrderEndpoint, po.ID)
}
