package main

import (
	"fmt"

	"github.com/stevenwilkin/swaps/binance"
)

var (
	b *binance.Binance
)

func main() {
	b = &binance.Binance{}

	for price := range b.Price() {
		fmt.Println(price)
	}
}
