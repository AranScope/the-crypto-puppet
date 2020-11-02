package coinbase

import (
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Account struct {
	ID             string          `json:"id"`
	Currency       string          `json:"currency"`
	Balance        decimal.Decimal `json:"balance"`
	Available      decimal.Decimal `json:"available"`
	Hold           decimal.Decimal `json:"hold"`
	ProfileID      string          `json:"profile_id"`
	TradingEnabled bool            `json:"trading_enabled"`
}

type AccountHistoryItem struct {
	ID        string          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Amount    decimal.Decimal `json:"amount"`
	Balance   decimal.Decimal `json:"balanace"`
	Type      string          `json:"type"`

	// If an entry is the result of a trade (match, fee), the details field will contain additional information about the trade.
	Details *struct {
		OrderID   string `json:"order_id"`
		TradeID   string `json:"trade_id"`
		ProductID string `json:"product_id"`
	} `json:"details"`
}

type Hold struct {
	ID        string          `json:"id"`
	AccountID string          `json:"account_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Amount    decimal.Decimal `json:"amount"`
	Type      string          `json:"type"`
	Ref       string          `json:"ref"`
}

func (c *Client) ListAccounts() ([]*Account, error) {
	var accounts []*Account
	err := c.Request(http.MethodGet, "/accounts", nil).DecodeResponse(&accounts)
	return accounts, err
}

func (c *Client) GetAccount(id string) (*Account, error) {
	var account *Account
	err := c.Request(http.MethodGet, "/accounts/"+id, nil).DecodeResponse(&account)
	return account, err
}

func (c *Client) GetHolds(accountID string) ([]*Hold, error) {
	var holds []*Hold
	err := c.Request(http.MethodGet, "/accounts/"+accountID+"/holds", nil).DecodeResponse(&holds)
	return holds, err
}

//func (c *Client) GetAccountHistory(id string) (chan<- *AccountHistoryItem, error)
