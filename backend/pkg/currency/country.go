package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"expensio-backend/internal/config"
)

type RestCountriesResponse struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Cca2 string `json:"cca2"` // Country code
}

// GetCountryCurrency fetches the default currency for a country using RestCountries API
func GetCountryCurrency(countryCode string, cfg *config.Config) (string, error) {
	url := fmt.Sprintf("%s/alpha/%s", cfg.ExternalAPIs.RestCountriesAPIURL, countryCode)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch country data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("rest countries API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var countries []RestCountriesResponse
	if err := json.Unmarshal(body, &countries); err != nil {
		return "", fmt.Errorf("failed to parse country response: %w", err)
	}

	if len(countries) == 0 {
		return "", fmt.Errorf("country not found")
	}

	country := countries[0]

	// Get the first currency (most countries have one primary currency)
	for currencyCode := range country.Currencies {
		return currencyCode, nil
	}

	return "", fmt.Errorf("no currency found for country")
}

// GetCountryCurrencyByName fetches currency by country name
func GetCountryCurrencyByName(countryName string, cfg *config.Config) (string, error) {
	url := fmt.Sprintf("%s/name/%s", cfg.ExternalAPIs.RestCountriesAPIURL, countryName)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch country data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("rest countries API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var countries []RestCountriesResponse
	if err := json.Unmarshal(body, &countries); err != nil {
		return "", fmt.Errorf("failed to parse country response: %w", err)
	}

	if len(countries) == 0 {
		return "", fmt.Errorf("country not found")
	}

	country := countries[0]

	// Get the first currency
	for currencyCode := range country.Currencies {
		return currencyCode, nil
	}

	return "", fmt.Errorf("no currency found for country")
}
