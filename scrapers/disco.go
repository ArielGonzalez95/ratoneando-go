package scrapers

import (
	"ratoneando/cores/catalog"
	"ratoneando/products"
)

func Disco(query string) ([]products.Schema, error) {
	return catalog.Core(catalog.CoreProps{
		Query:   query,
		BaseUrl: "https://www.disco.com.ar",
		Source:  "disco",
	})
}
