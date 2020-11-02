package coinbase

import (
	"github.com/shopspring/decimal"
	"net/http"
)

type Fees struct {
	MakerFeeRate decimal.Decimal `json:"maker_fee_rate"`
	TakerFeeRate decimal.Decimal `json:"taker_fee_rate"`
	USDVolume    decimal.Decimal `json:"usd_volume"`
}

func (c *Client) GetFees() (*Fees, error) {
	var fees *Fees
	err := c.Request(http.MethodGet, "/fees", nil).DecodeResponse(&fees)
	return fees, err
}
