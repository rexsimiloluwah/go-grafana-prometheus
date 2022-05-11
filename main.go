package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rexsimiloluwah/go-prometheus-grafana-demo/middlewares"
)

type QuoteResponse struct {
	ID           string   `json:"_id"`
	Tags         []string `json:"tags"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	AuthorSlug   string   `json:"authorSlug"`
	Length       int      `json:"length"`
	DateAdded    string   `json:"dateAdded"`
	DateModified string   `json:"dateModified"`
}

const (
	RANDOM_QUOTE_API_ENDPOINT = "https://api.quotable.io/random"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	e := echo.New()
	middlewares.RegisterPrometheusMetrics()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "\n[${status}]: ${method} ${uri}",
	}))

	// Endpoint for serving the scraped prometheus metrics
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/quote", RandomQuoteHandler, middlewares.RecordRequestLatency)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", PORT)))
}

func RandomQuoteHandler(c echo.Context) error {
	var quote QuoteResponse
	resp, err := http.Get(RANDOM_QUOTE_API_ENDPOINT)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	err = json.NewDecoder(resp.Body).Decode(&quote)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully generated random quote.",
		"data":    quote,
	})
}
