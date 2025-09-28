package money

import (
	"fmt"
	"strings"
)

type Currency string

const (
	PLN Currency = "PLN"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

var AllCurrencies = []Currency{
	USD, EUR, PLN,
}

func (c Currency) String() string {
	return string(c)
}

func (c Currency) IsValid() bool {
	for _, currency := range AllCurrencies {
		if currency == c {
			return true
		}
	}
	return false
}

func ParseCurrency(s string) (Currency, error) {
	currency := Currency(strings.ToUpper(s))
	if !currency.IsValid() {
		return "", fmt.Errorf("invalid currency: %s", s)
	}
	return currency, nil
}

func GetCurrencySymbol(currency Currency) string {
	symbols := map[Currency]string{
		USD: "$",
		EUR: "€",
		PLN: "zł",
	}

	if symbol, exists := symbols[currency]; exists {
		return symbol
	}
	return currency.String()
}
