package scrapers

import (
	"ratoneando/cores/catalog"
	"ratoneando/products"
)

func Vea(query string) ([]products.Schema, error) {
	return catalog.Core(catalog.CoreProps{
		Query:   query,
		BaseUrl: "https://www.vea.com.ar",
		Source:  "vea",
	})
}
