package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type model struct {
}

func initialModel() model {

	return model{}

}

func (m model) Init() tea.Cmd {

	return nil

}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil

}

func (m model) View() string {

	return pterm.DefaultCenter.Sprint(
		pterm.DefaultBigText.
			WithLetters(putils.LettersFromString("Hello World!")),
	)

}

func main() {

	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {

		log.Fatal(err)

	}

}
