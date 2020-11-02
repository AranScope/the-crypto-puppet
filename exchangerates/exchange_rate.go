package exchangerates

import "github.com/shopspring/decimal"

type Rates struct {
	Rates map[string]decimal.Decimal `json:"rates"`
}

func (c *Client) GetExchangeRates(baseCurrency string) (*Rates, error) {
	var rates *Rates
	err := c.Request("GET", "/latest?base="+baseCurrency, nil).DecodeResponse(&rates)
	return rates, err
}
