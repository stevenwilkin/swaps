package main

import (
	"fmt"
	"os"

	"github.com/stevenwilkin/swaps/binance"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	b      *binance.Binance
	margin = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type spotMsg float64

type model struct {
	spot float64
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
	}

	return m, nil
}

func (m model) View() string {
	spot := fmt.Sprintf("Spot: %7.2f", m.spot)

	return margin.Render(fmt.Sprintf(
		"%s", spot))
}

func main() {
	m := model{}
	p := tea.NewProgram(m, tea.WithAltScreen())

	b = &binance.Binance{}

	go func() {
		for price := range b.Price() {
			p.Send(spotMsg(price))
		}
	}()

	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
