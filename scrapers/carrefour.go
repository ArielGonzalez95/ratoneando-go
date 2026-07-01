package scrapers

import (
	"ratoneando/cores/catalog"
	"ratoneando/products"
)

func Carrefour(query string) ([]products.Schema, error) {
	return catalog.Core(catalog.CoreProps{
		Query:   query,
		BaseUrl: "https://www.carrefour.com.ar",
		Source:  "carrefour",
	})
}
