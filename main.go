package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stevenwilkin/carry/feed"
	"github.com/stevenwilkin/swaps/binance"
	"github.com/stevenwilkin/swaps/bybit"
	"github.com/stevenwilkin/swaps/deribit"
)

var (
	h  = feed.NewHandler()
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

func displayFailing() {
	fmt.Println("\033[2J\033[H\033[?25l") // clear screen, move cursor to top of screen, hide cursor
	fmt.Println("Feed failing...")
}

func exitFailed() {
	fmt.Println("Feed failed")
	os.Exit(1)
}

func main() {
	h.Add(feed.NewFeed(b.Price, func(x float64) {
		spot = x
	}))

	h.Add(feed.NewFeed(by.Price, func(x float64) {
		sBybit = x
	}))

	h.Add(feed.NewFeed(d.Price, func(x float64) {
		sDeribit = x
	}))

	t := time.NewTicker(100 * time.Millisecond)

	for {
		if h.Failed() {
			exitFailed()
		} else if h.Failing() {
			displayFailing()
		} else {
			display()
		}

		<-t.C
	}
}
