package coinbase

import (
	"net/http"
)

func (c *Client) GetPortfolios() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	err := c.Request(http.MethodGet, "/currencies", nil).DecodeResponse(&result)
	return result, err
}
