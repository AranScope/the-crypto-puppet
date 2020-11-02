package functions

import (
	"encoding/json"
	"fmt"
	"github.com/AranScope/the-crypto-puppet/coinbase"
	"github.com/shopspring/decimal"
	"math/rand"
	"reflect"
)

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
}

func Prettify(v interface{}) string {
	bytes, _ := json.MarshalIndent(v, "", "  ")
	return reflect.TypeOf(v).String() + "\n===============\n\n" + string(bytes) + "\n\n===============\n"
}

func filterOnlyAccountsWithTradingEnabled(accounts []*coinbase.Account) []*coinbase.Account {
	var filteredAccounts []*coinbase.Account
	for _, account := range accounts {
		if account.TradingEnabled {
			filteredAccounts = append(filteredAccounts, account)
		}
	}
	return filteredAccounts
}

func filterOnlyAccountsWithPositiveAvailableBalances(accounts []*coinbase.Account) []*coinbase.Account {
	var filteredAccounts []*coinbase.Account
	for _, account := range accounts {
		if account.Available.GreaterThan(decimal.Zero) {
			filteredAccounts = append(filteredAccounts, account)
		}
	}
	return filteredAccounts
}

func randomAccount(accounts []*coinbase.Account) *coinbase.Account {
	return accounts[rand.Intn(len(accounts))]
}

func randomProduct(products []*coinbase.Product) *coinbase.Product {
	return products[rand.Intn(len(products))]
}

func randomSignificantAmount(account *coinbase.Account) decimal.Decimal {
	// Choose a random amount between 25% and 100%
	return account.Available.Sub(account.Available.Mul(decimal.RequireFromString(fmt.Sprintf("0.%d", rand.Intn(75)))))
}