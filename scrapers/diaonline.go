package scrapers

import (
	"ratoneando/cores/catalog"
	"ratoneando/products"
)

func DiaOnline(query string) ([]products.Schema, error) {
	return catalog.Core(catalog.CoreProps{
		Query:   query,
		BaseUrl: "https://diaonline.supermercadosdia.com.ar",
		Source:  "diaonline",
	})
}
