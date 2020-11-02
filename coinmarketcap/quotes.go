package coinmarketcap

import (
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
)

type Quote struct {
	Data map[string]struct {
		Quote struct {
			USD struct {
				Price            decimal.Decimal `json:"price"`
				PercentChange24h decimal.Decimal `json:"percent_change_24h"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

func (c *Client) GetQuote(symbols ...string) (*Quote, error) {
	var quote *Quote
	err := c.Request(http.MethodGet, "/v1/cryptocurrency/quotes/latest?symbol="+strings.Join(symbols, ","), nil).DecodeResponse(&quote)
	return quote, err
}
