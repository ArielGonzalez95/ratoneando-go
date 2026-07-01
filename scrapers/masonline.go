package scrapers

import (
	"ratoneando/cores/catalog"
	"ratoneando/products"
)

func MasOnline(query string) ([]products.Schema, error) {
	return catalog.Core(catalog.CoreProps{
		Query:   query,
		BaseUrl: "https://www.masonline.com.ar",
		Source:  "masonline",
	})
}
