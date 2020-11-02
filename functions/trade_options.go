package functions

import (
	"fmt"
	"github.com/AranScope/the-crypto-puppet/coinbase"
	"github.com/shopspring/decimal"
	"math/rand"
	"os"
	"strings"
	"time"
)

var tradingIntroStrings = []string{
	"It's trading time! ⏰",
	"Let's get cooking with another trade! 🍲",
	"I'm feeling another trade! 🔜",
	"...Yep, another trade! 🥱",
	"I don't think we've had enough trades already! 💨",
	"I live to trade! 🧬",
}

type TradeOption struct {
	FromSymbol string
	ToSymbol   string
	Amount     decimal.Decimal
}

func SendTradeOptionsPoll() (string, error) {
	tradeOptions, err := generateRandomTradeOptions()
	if err != nil {
		return "", err
	}

	body := generateTweetBody(tradeOptions)
	//sendTweet(body)
	return body, nil
}

func generateTweetBody(options []*TradeOption) string {
	intro := fmt.Sprintf("🤖 %s\n\n⚠ Voting closes in 1 hour.\n\n", tradingIntroStrings[rand.Intn(len(tradingIntroStrings))])

	var tradeLines []string
	for i, tradeOption := range options {
		tradeLines = append(tradeLines, fmt.Sprintf("%d: %s %s ➡ %s", i+1, tradeOption.Amount.StringFixed(2), tradeOption.FromSymbol, tradeOption.ToSymbol))
	}
	tradeLines = append(tradeLines, fmt.Sprintf("%d: HODL (Do nothing)", len(options)+1))

	body := intro + strings.Join(tradeLines, "\n")
	return body
}

func generateRandomTradeOptions() ([]*TradeOption, error) {
	rand.Seed(time.Now().Unix())
	client := coinbase.NewSandboxClient(
		os.Getenv("COINBASE-ACCESS-KEY"),
		os.Getenv("COINBASE-ACCESS-SECRET"),
		os.Getenv("COINBASE-ACCESS-PASSPHRASE"),
	)

	accounts, err := client.ListAccounts()
	if err != nil {
		return nil, err
	}

	products, err := client.GetProducts()
	if err != nil {
		return nil, err
	}
	productsByBaseCurrency := map[string][]*coinbase.Product{}
	for _, product := range products {
		productsByBaseCurrency[product.BaseCurrency] = append(productsByBaseCurrency[product.BaseCurrency], product)
	}

	tradingAccounts := filterOnlyAccountsWithTradingEnabled(accounts)

	// We can only trade 'from' accounts with positive balances
	accountsWithPositiveBalances := filterOnlyAccountsWithPositiveAvailableBalances(tradingAccounts)

	var accountsWithTradingProducts []*coinbase.Account
	for _, account := range accountsWithPositiveBalances {
		if _, ok := productsByBaseCurrency[account.Currency]; ok {
			accountsWithTradingProducts = append(accountsWithTradingProducts, account)
		}
	}

	numTrades := min(len(accountsWithPositiveBalances), 3)
	var tradeOptions []*TradeOption

	for i := 0; i < numTrades; i++ {
		fromAccount := randomAccount(accountsWithTradingProducts)
		toSymbol := randomProduct(productsByBaseCurrency[fromAccount.Currency]).QuoteCurrency
		tradeOptions = append(tradeOptions, &TradeOption{
			FromSymbol: fromAccount.Currency,
			ToSymbol:   toSymbol,
			Amount:     randomSignificantAmount(fromAccount),
		})
	}

	return tradeOptions, nil
}
