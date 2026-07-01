package catalog

import (
	"encoding/json"
	"net/url"

	"ratoneando/cores/api"
	"ratoneando/products"
)

type CoreProps struct {
	Query   string
	BaseUrl string
	Source  string
}

type responseProduct struct {
	ProductId   string   `json:"productId"`
	ProductName string   `json:"productName"`
	Link        string   `json:"link"`
	ProductData []string `json:"ProductData"`
	Items       []struct {
		Images []struct {
			ImageUrl string `json:"imageUrl"`
		} `json:"images"`
		Sellers []struct {
			CommertialOffer struct {
				Price             float64 `json:"Price"`
				ListPrice         float64 `json:"ListPrice"`
				AvailableQuantity int     `json:"AvailableQuantity"`
				IsAvailable       bool    `json:"IsAvailable"`
			} `json:"commertialOffer"`
		} `json:"sellers"`
	} `json:"items"`
}

type productMeta struct {
	MeasurementUnit string  `json:"MeasurementUnit"`
	UnitMultiplier  float64 `json:"UnitMultiplier"`
}

type rawProduct struct {
	responseProduct
	productMeta
}

func Core(props CoreProps) ([]products.Schema, error) {
	return api.Core(api.CoreProps[[]responseProduct, rawProduct]{
		Query:   props.Query,
		BaseUrl: props.BaseUrl,
		SearchPattern: func(q string) string {
			return "/api/catalog_system/pub/products/search/?ft=" + url.QueryEscape(q) + "&_from=0&_to=19"
		},
		Source: props.Source,
		Normalizer: func(response []responseProduct) []rawProduct {
			result := make([]rawProduct, 0, len(response))
			for _, p := range response {
				var meta productMeta
				if len(p.ProductData) > 0 {
					json.Unmarshal([]byte(p.ProductData[0]), &meta)
				}
				result = append(result, rawProduct{responseProduct: p, productMeta: meta})
			}
			return result
		},
		Extractor: func(p rawProduct) products.ExtendedSchema {
			if len(p.Items) == 0 || len(p.Items[0].Sellers) == 0 {
				return products.ExtendedSchema{}
			}
			offer := p.Items[0].Sellers[0].CommertialOffer
			imageUrl := ""
			if len(p.Items[0].Images) > 0 {
				imageUrl = p.Items[0].Images[0].ImageUrl
			}
			return products.ExtendedSchema{
				ID:          p.ProductId,
				Source:      props.Source,
				Name:        p.ProductName,
				Link:        p.Link,
				Image:       imageUrl,
				Unavailable: !offer.IsAvailable,
				Price:       offer.Price,
				ListPrice:   offer.ListPrice,
				Unit:        p.MeasurementUnit,
				UnitFactor:  p.UnitMultiplier,
			}
		},
	})
}
