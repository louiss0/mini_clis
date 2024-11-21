package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mini-clis/form"
)

func main() {

	if _, err := tea.NewProgram(form.InitialModel(), tea.WithAltScreen()).Run(); err != nil {

		log.Fatal(err)

	}

}
