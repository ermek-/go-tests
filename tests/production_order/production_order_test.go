package production_order

import (
	"fmt"
	"net/http"
	"testing"

	"go-api-tests/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProductionOrder(t *testing.T) {
	c := helpers.TestClient(t)

	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
	req := CreateRequest{Number: &number}

	resp, po := Create(t, c, Endpoints, req)

	require.Equal(t, http.StatusCreated, resp.StatusCode, "expected 201 Created")
	require.Equal(t, number, po.Number, "expected production order number to match")
	require.NotZero(t, po.ID, "server must return non-zero id")
	require.Equal(t, "Создан", po.Status, "expected production order status to match")

	getResp, got := Get(t, c, Endpoint, po.ID)
	require.Equal(t, http.StatusOK, getResp.StatusCode, "expected 200 OK on GET by id")
	assert.Equal(t, po.ID, got.ID, "GET should return the same ID")
	assert.Equal(t, number, got.Number, "GET should return the same Number")
	assert.Equal(t, "Создан", got.Status, "status should match after create")

	t.Cleanup(func() { Delete(t, c, Endpoint, po.ID) })
}
