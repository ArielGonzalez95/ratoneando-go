package controllers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"ratoneando/products"
	"ratoneando/scrapers"
)

type ListaRequest struct {
	Items []string `json:"items"`
}

type ItemResult struct {
	Query   string           `json:"query"`
	Best    *products.Schema `json:"best"`
	Options []products.Schema `json:"options"`
}

type ListaResponse struct {
	Results   []ItemResult         `json:"results"`
	Canasta   map[string][]string  `json:"canasta"`
	Timestamp time.Time            `json:"timestamp"`
}

var allScrapers = []func(string) ([]products.Schema, error){
	scrapers.Carrefour,
	scrapers.Coto,
	scrapers.DiaOnline,
	scrapers.Disco,
	scrapers.Farmacity,
	scrapers.Jumbo,
	scrapers.MasOnline,
	scrapers.MercadoLibre,
	scrapers.Vea,
}

func Lista(c *gin.Context) {
	var req ListaRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere al menos un item."})
		return
	}

	var wg sync.WaitGroup
	results := make([]ItemResult, len(req.Items))

	for i, item := range req.Items {
		wg.Add(1)
		go func(idx int, raw string) {
			defer wg.Done()

			query := strings.ToLower(strings.TrimSpace(raw))
			if query == "" {
				results[idx] = ItemResult{Query: raw}
				return
			}

			var scraperWg sync.WaitGroup
			var mu sync.Mutex
			var allProds []products.Schema

			for _, scraper := range allScrapers {
				scraperWg.Add(1)
				go func(s func(string) ([]products.Schema, error)) {
					defer scraperWg.Done()
					p, err := s(query)
					if err != nil || len(p) == 0 {
						return
					}
					mu.Lock()
					allProds = append(allProds, p...)
					mu.Unlock()
				}(scraper)
			}
			scraperWg.Wait()

			filtered := products.Fuzzy(allProds, query)
			sorted := products.Sort(filtered)

			var best *products.Schema
			if len(sorted) > 0 {
				b := sorted[0]
				best = &b
			}

			opts := sorted
			if len(opts) > 5 {
				opts = opts[:5]
			}

			results[idx] = ItemResult{
				Query:   raw,
				Best:    best,
				Options: opts,
			}
		}(i, item)
	}

	wg.Wait()

	// Canasta óptima: agrupar por supermarket donde cada item sale más barato
	canasta := make(map[string][]string)
	for _, r := range results {
		if r.Best != nil {
			canasta[r.Best.Source] = append(canasta[r.Best.Source], r.Query)
		}
	}

	c.JSON(http.StatusOK, ListaResponse{
		Results:   results,
		Canasta:   canasta,
		Timestamp: time.Now(),
	})
}
