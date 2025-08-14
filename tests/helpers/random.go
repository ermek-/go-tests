package helpers

import "github.com/brianvoe/gofakeit/v6"

func init() { gofakeit.Seed(0) }

func RandomNumber() int {
	return gofakeit.Number(1, 10000)
}
