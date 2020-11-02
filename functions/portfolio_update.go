package functions

import (
	"fmt"
	"github.com/AranScope/the-crypto-puppet/coinbase"
	"github.com/AranScope/the-crypto-puppet/coinmarketcap"
	"github.com/AranScope/the-crypto-puppet/exchangerates"
	"github.com/shopspring/decimal"
	"os"
	"sort"
)

func SendPortfolioUpdate() error {
	balances, err := GetPortfolioBalances()
	if err != nil {
		return err
	}

	tweet := "🤖 Hey hey 👋 Here's a breakdown of our portfolio top 5:\n"

	for i := 0; i < min(len(balances), 5); i++ {
		balance := balances[i]
		tweet += fmt.Sprintf("- %s %s", balance.Units.StringFixed(2), balance.Currency)

		if balance.Currency != "USD" {
			tweet += fmt.Sprintf(" (~$%s)", balance.USDValue.StringFixed(2))
		}
		tweet += "\n"
	}

	totalValue := decimal.Zero
	for _, bal := range balances {
		totalValue = totalValue.Add(bal.USDValue)
	}

	tweet += "\nTotal portfolio value: $" + totalValue.StringFixed(2)

	fmt.Println(tweet)

	return nil
}

type Balance struct {
	Currency string
	Units    decimal.Decimal
	USDValue decimal.Decimal
}

func GetPortfolioBalances() ([]Balance, error) {
	var balances []Balance

	client := coinbase.NewSandboxClient(
		os.Getenv("COINBASE-ACCESS-KEY"),
		os.Getenv("COINBASE-ACCESS-SECRET"),
		os.Getenv("COINBASE-ACCESS-PASSPHRASE"),
	)

	accounts, err := client.ListAccounts()
	if err != nil {
		return nil, err
	}

	var accountsWithPositiveBalances []*coinbase.Account
	for _, account := range accounts {
		if account.Balance.GreaterThan(decimal.Zero) {
			accountsWithPositiveBalances = append(accountsWithPositiveBalances, account)
		}
	}

	exClient := exchangerates.NewClient()
	fiatExchangeRates, err := exClient.GetExchangeRates("USD")
	if err != nil {
		return nil, err
	}

	var symbols []string
	for _, account := range accountsWithPositiveBalances {
		if isFiat(account.Currency) {
			if account.Currency == "USD" {
				balances = append(balances, Balance{
					Currency: account.Currency,
					USDValue: account.Balance,
					Units:    account.Balance,
				})
			} else {
				balances = append(balances, Balance{
					Currency: account.Currency,
					Units:    account.Balance,
					USDValue: account.Balance.Div(fiatExchangeRates.Rates[account.Currency]),
				})
			}
		} else {
			symbols = append(symbols, account.Currency)
		}
	}

	cmcClient := coinmarketcap.NewClient(os.Getenv("COINMARKETCAP_ACCESS_SECRET"))
	quote, err := cmcClient.GetQuote(symbols...)
	if err != nil {
		return nil, err
	}

	for _, account := range accountsWithPositiveBalances {
		if !isFiat(account.Currency) {
			balances = append(balances, Balance{
				Currency: account.Currency,
				Units:    account.Balance,
				USDValue: account.Balance.Mul(quote.Data[account.Currency].Quote.USD.Price),
			})
		}
	}

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].USDValue.LessThan(balances[j].USDValue)
	})

	return balances, nil
}

func isFiat(currency string) bool {
	return currency == "USD" || currency == "EUR" || currency == "GBP"
}
