package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"ratoneando/products"
	"ratoneando/unit"
	"ratoneando/utils/logger"
)

var httpClient = &http.Client{}

func Core[ResponseStructure any, RawProduct any](props CoreProps[ResponseStructure, RawProduct]) ([]products.Schema, error) {
	escapedQuery := url.PathEscape(props.Query)
	searchUrl := props.BaseUrl + props.SearchPattern(escapedQuery)

	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		logger.LogError("Failed to create request: " + escapedQuery + "@" + props.Source)
		return nil, fmt.Errorf(props.Source)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "es-AR,es;q=0.9")

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.LogError("Failed to fetch the URL: " + escapedQuery + "@" + props.Source)
		return nil, fmt.Errorf(props.Source)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.LogError(fmt.Sprintf("HTTP %d from %s for query %s", resp.StatusCode, props.Source, escapedQuery))
		return nil, fmt.Errorf(props.Source)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError("Failed to read the response body: " + escapedQuery + "@" + props.Source)
		return nil, fmt.Errorf(props.Source)
	}

	var responseStructure ResponseStructure
	err = json.Unmarshal(body, &responseStructure)
	if err != nil {
		logger.LogError("Failed to unmarshal the response body: " + escapedQuery + "@" + props.Source)
		return nil, fmt.Errorf(props.Source)
	}

	var errorCheck struct {
		Errors []struct {
			Message    string `json:"message"`
			Extensions struct {
				Code string `json:"code"`
			} `json:"extensions"`
			Name string `json:"name"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &errorCheck); err == nil && len(errorCheck.Errors) > 0 {
		logger.LogError("API returned error: " + errorCheck.Errors[0].Message + " for " + escapedQuery + "@" + props.Source)
		return nil, fmt.Errorf(props.Source)
	}

	normalizedProducts := props.Normalizer(responseStructure)

	products := make([]products.Schema, len(normalizedProducts))
	for i, product := range normalizedProducts {
		extractedProduct := props.Extractor(product)
		if extractedProduct.Unavailable {
			continue
		}
		products[i] = unit.CalculateUnitInfo(extractedProduct)
	}

	return products, nil
}
