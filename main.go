package main

import "github.com/AranScope/the-crypto-puppet/functions"

func main() {
	err := functions.SendPortfolioUpdate()
	if err != nil {
		panic(err)
	}
}
