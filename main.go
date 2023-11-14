package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/swaps/binance"
	"github.com/stevenwilkin/swaps/bybit"
	"github.com/stevenwilkin/swaps/deribit"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"
)

var (
	b      = &binance.Binance{}
	by     = bybit.NewBybitFromEnv()
	d      = &deribit.Deribit{}
	margin = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type spotMsg float64
type bybitMsg float64
type deribitMsg float64

type model struct {
	spot    float64
	bybit   float64
	deribit float64
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case spotMsg:
		m.spot = float64(msg)
	case bybitMsg:
		m.bybit = float64(msg)
	case deribitMsg:
		m.deribit = float64(msg)
	}

	return m, nil
}

func delta(spot, perp float64) string {
	if perp == 0 || spot == 0 {
		return ""
	}

	return fmt.Sprintf("%6.2f", perp-spot)
}

func (m model) View() string {
	spot := fmt.Sprintf("Spot:    %8.2f", m.spot)
	bybit := fmt.Sprintf("Bybit:   %8.2f %s", m.bybit, delta(m.spot, m.bybit))
	deribit := fmt.Sprintf("Deribit: %8.2f %s", m.deribit, delta(m.spot, m.deribit))

	return margin.Render(fmt.Sprintf(
		"%s\n%s\n%s", spot, bybit, deribit))
}

func main() {
	m := model{}
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for price := range b.Price() {
			p.Send(spotMsg(price))
		}
		os.Exit(1)
	}()

	go func() {
		for price := range by.Price() {
			p.Send(bybitMsg(price))
		}
		os.Exit(1)
	}()

	go func() {
		for price := range d.Price() {
			p.Send(deribitMsg(price))
		}
		os.Exit(1)
	}()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
