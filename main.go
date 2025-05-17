package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stevenwilkin/swaps/binance"
	"github.com/stevenwilkin/swaps/bybit"
	"github.com/stevenwilkin/swaps/deribit"
)

var (
	b  = &binance.Binance{}
	by = &bybit.Bybit{}
	d  = &deribit.Deribit{}

	spot     float64
	sBybit   float64
	sDeribit float64
)

func delta(spot, perp float64) string {
	if perp == 0 || spot == 0 {
		return ""
	}

	return fmt.Sprintf("%6.2f", perp-spot)
}

func display() {
	fmt.Println("\033[2J\033[H\033[?25l") // clear screen, move cursor to top of screen, hide cursor

	fmt.Printf("  Spot:    %9.2f\n", spot)
	fmt.Printf("  Bybit:   %9.2f %s\n", sBybit, delta(spot, sBybit))
	fmt.Printf("  Deribit: %9.2f %s\n", sDeribit, delta(spot, sDeribit))
}

func process[T any](f func() chan T, p func(x T)) {
	go func() {
		for result := range f() {
			p(result)
		}
		os.Exit(1)
	}()
}

func main() {
	process(b.Price, func(x float64) {
		spot = x
	})

	process(by.Price, func(x float64) {
		sBybit = x
	})

	process(d.Price, func(x float64) {
		sDeribit = x
	})

	t := time.NewTicker(100 * time.Millisecond)

	for {
		display()
		<-t.C
	}
}
