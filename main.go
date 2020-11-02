package main

import (
	"fmt"
	"github.com/AranScope/the-crypto-puppet/functions"
)

func main() {
	err := functions.SendPortfolioUpdate()
	if err != nil {
		panic(err)
	}

	t, err := functions.SendTradeOptionsPoll()
	if err != nil {
		panic(err)
	}

	fmt.Println(functions.Prettify(t))
}
