package tests

import (
	"os"
	"testing"

	"go-api-tests/internal/api"
	"go-api-tests/internal/env"
)

var (
	client                   *api.Client
	nomenclatureEndpoint     string
	nomenclaturesEndpoint    string
	productionOrderEndpoint  string
	productionOrdersEndpoint string
)

func TestMain(m *testing.M) {
	_ = env.LoadDotEnv(".env")
	baseURL := api.Env("BASE_URL", "http://api.dev.pkmt.tech")
	authEndpoint := api.Env("AUTH_ENDPOINT", "")
	username := api.Env("USERNAME", "")
	password := api.Env("PASSWORD", "")

	nomenclatureEndpoint = "/Nomenclature/v1/nomenclatures"
	nomenclaturesEndpoint = "/Nomenclature/v1/nomenclatures/"
	productionOrderEndpoint = "/ProductionOrder/v1/ProductionOrders"
	productionOrdersEndpoint = "/ProductionOrder/v1/ProductionOrders/"

	client = api.NewClient(baseURL)

	if err := client.Authenticate(authEndpoint, username, password); err != nil && authEndpoint != "" {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}
