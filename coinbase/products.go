package coinbase

import (
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Product struct {
	ID              string          `json:"id"`
	DisplayName     string          `json:"display_name"`
	BaseCurrency    string          `json:"base_currency"`
	QuoteCurrency   string          `json:"quote_currency"`
	BaseIncrement   decimal.Decimal `json:"base_increment"`
	QuoteIncrement  decimal.Decimal `json:"quote_increment"`
	BaseMinSize     decimal.Decimal `json:"base_min_size"`
	BaseMaxSize     decimal.Decimal `json:"base_max_size"`
	MinMarketFunds  decimal.Decimal `json:"min_market_funds"`
	MaxMarketFunds  decimal.Decimal `json:"max_market_funds"`
	Status          string          `json:"status"`
	StatusMessage   string          `json:"status_message"`
	CancelOnly      bool            `json:"cancel_only"`
	LimitOnly       bool            `json:"limit_only"`
	PostOnly        bool            `json:"post_only"`
	TradingDisabled bool            `json:"trading_disabled"`
}

type Ticker struct {
	TradeID string          `json:"trade_id"`
	Price   decimal.Decimal `json:"price"`
	Size    decimal.Decimal `json:"size"`
	Bid     decimal.Decimal `json:"bid"`
	Ask     decimal.Decimal `json:"ask"`
	Volume  decimal.Decimal `json:"volume"`
	Time    time.Time       `json:"time"`
}

func (c *Client) GetProducts() ([]*Product, error) {
	var products []*Product
	err := c.Request(http.MethodGet, "/products", nil).DecodeResponse(&products)
	return products, err
}

func (c *Client) GetProduct(ID string) (*Product, error) {
	var product *Product
	err := c.Request(http.MethodGet, "/products/"+ID, nil).DecodeResponse(&product)
	return product, err
}
func (c *Client) GetProductTicket(ID string) (*Ticker, error) {
	var ticker *Ticker
	err := c.Request(http.MethodGet, "/products/"+ID+"/ticket", nil).DecodeResponse(&ticker)
	return ticker, err
}
