package coinbase

import (
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type OrderRequest struct {
	ClientOID string          `json:"client_oid,omitempty"` // optional uuid, this will be broadcast to the public feed and can be used to identify your orders
	Type      string          `json:"type,omitempty"`       // optional 'limit' or 'market', default is limit
	Side      string          `json:"side"`                 // 'buy' or 'sell'
	ProductID string          `json:"product_id"`           // must match a valid product from the /products list
	STP       string          `json:"stp,omitempty"`        // optional
	Stop      string          `json:"stop,omitempty"`       // optional 'loss' or 'entry', requires 'stop_price' to be defined
	StopPrice decimal.Decimal `json:"stop_price,omitempty"` // optional, only if `stop` is defined

	// shared order parameters
	Size decimal.Decimal `json:"size,omitempty"` // required if 'type' is limit, one of 'size' or 'funds' must be set for market orders

	// limit order parameters, requires 'type' to be 'limit'
	Price       decimal.Decimal `json:"price,omitempty"`         // required if 'type' is limit
	TimeInForce string          `json:"time_in_force,omitempty"` // optional, 'GTC', 'GTT', 'IOC', 'FOK', default is 'GTC'
	CancelAfter time.Time       `json:"cancel_after,omitempty"`  // optional, requires 'time_in_force' to be 'GTT'
	PostOnly    bool            `json:"post_only,omitempty"`     // optional, invalid when 'time_in_force' is 'IOC' or 'FOK'

	// market order parameters, requires 'type' to be 'market'
	Funds decimal.Decimal `json:"funds"` // one of 'size' or 'funds' must be set for market orders
}

type Order struct {
	ID            string          `json:"id"`
	Price         decimal.Decimal `json:"price"`
	Size          decimal.Decimal `json:"size"`
	ProductID     string          `json:"product_id"`
	Side          string          `json:"side"`
	STP           string          `json:"stp"`
	Type          string          `json:"type"`
	TimeInForce   string          `json:"time_in_force"`
	PostOnly      bool            `json:"post_only"`
	CreatedAt     time.Time       `json:"created_at"`
	DoneAt        time.Time       `json:"done_at"`
	FillFees      decimal.Decimal `json:"fill_fees"`
	FilledSize    decimal.Decimal `json:"filled_size"`
	ExecutedValue decimal.Decimal `json:"executed_value"`
	Status        string          `json:"status"`
	Settled       bool            `json:"settled"`
}

func (c *Client) PlaceOrder(req *OrderRequest) (*Order, error) {
	var rsp *Order
	err := c.Request(http.MethodGet, "/orders", req).DecodeResponse(&rsp)
	return rsp, err
}

func (c *Client) GetOrder(id string) (*Order, error) {
	var order *Order
	err := c.Request(http.MethodGet, "/orders/"+id, nil).DecodeResponse(&order)
	return order, err
}

func (c *Client) GetOrderByClientOID(oid string) (*Order, error) {
	var order *Order
	err := c.Request(http.MethodGet, "/orders/client:"+oid, nil).DecodeResponse(&order)
	return order, err
}

func (c *Client) CancelOrder(id string) error {
	return c.Request(http.MethodDelete, "/orders/"+id, nil).DecodeResponse(nil)
}

func (c *Client) CancelOrderByClientOID(oid string) error {
	return c.Request(http.MethodDelete, "/orders/client:"+oid, nil).DecodeResponse(nil)
}
