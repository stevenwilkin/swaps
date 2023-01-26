package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/swaps/binance"
	"github.com/stevenwilkin/swaps/bybit"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	b      = &binance.Binance{}
	by     = &bybit.Bybit{}
	margin = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type spotMsg float64
type bybitMsg float64

type model struct {
	spot  float64
	bybit float64
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
	}

	return m, nil
}

func (m model) View() string {
	spot := fmt.Sprintf("Spot:  %7.2f", m.spot)
	bybit := fmt.Sprintf("Bybit: %7.2f", m.bybit)

	return margin.Render(fmt.Sprintf(
		"%s\n%s", spot, bybit))
}

func main() {
	m := model{}
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for price := range b.Price() {
			p.Send(spotMsg(price))
		}
	}()

	go func() {
		for price := range by.Price() {
			p.Send(bybitMsg(price))
		}
	}()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
