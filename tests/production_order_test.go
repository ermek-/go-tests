package tests

import (
	"fmt"
	"testing"

	"go-api-tests/tests/helpers"
)

func TestCreateProductionOrder(t *testing.T) {
	number := fmt.Sprintf("TEST-%d", helpers.RandomNumber())
	po := helpers.CreateProductionOrder(t, client, productionOrdersEndpoint, helpers.ProductionOrder{Number: &number})
	helpers.DeleteProductionOrder(t, client, productionOrderEndpoint, po.ID)
}
