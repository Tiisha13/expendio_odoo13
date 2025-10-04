package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"expensio-backend/internal/config"
	"expensio-backend/pkg/cache"
)

type ExchangeRateResponse struct {
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// GetExchangeRate fetches exchange rate from API or cache
func GetExchangeRate(from, to string, cfg *config.Config) (float64, error) {
	// If same currency, return 1.0
	if from == to {
		return 1.0, nil
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("exchange_rate:%s:%s", from, to)
	var rate float64
	err := cache.Get(cacheKey, &rate)
	if err == nil {
		return rate, nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/%s", cfg.ExternalAPIs.ExchangeRateAPIURL, from)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("exchange rate API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var exchangeRateResp ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRateResp); err != nil {
		return 0, fmt.Errorf("failed to parse exchange rate response: %w", err)
	}

	// Get the rate for target currency
	rate, exists := exchangeRateResp.Rates[to]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for currency: %s", to)
	}

	// Cache the result
	_ = cache.Set(cacheKey, rate, cfg.Cache.CurrencyRateTTL)

	return rate, nil
}

// ConvertCurrency converts amount from one currency to another
func ConvertCurrency(amount float64, from, to string, cfg *config.Config) (float64, float64, error) {
	rate, err := GetExchangeRate(from, to, cfg)
	if err != nil {
		return 0, 0, err
	}

	convertedAmount := amount * rate
	return convertedAmount, rate, nil
}

// GetAllRates fetches all exchange rates for a base currency
func GetAllRates(baseCurrency string, cfg *config.Config) (map[string]float64, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("exchange_rates:%s", baseCurrency)
	var rates map[string]float64
	err := cache.Get(cacheKey, &rates)
	if err == nil {
		return rates, nil
	}

	// Fetch from API
	url := fmt.Sprintf("%s/%s", cfg.ExternalAPIs.ExchangeRateAPIURL, baseCurrency)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("exchange rate API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var exchangeRateResp ExchangeRateResponse
	if err := json.Unmarshal(body, &exchangeRateResp); err != nil {
		return nil, fmt.Errorf("failed to parse exchange rate response: %w", err)
	}

	// Cache the result
	_ = cache.Set(cacheKey, exchangeRateResp.Rates, cfg.Cache.CurrencyRateTTL)

	return exchangeRateResp.Rates, nil
}
