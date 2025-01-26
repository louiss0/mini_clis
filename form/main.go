package form

import (
	"fmt"
	"time"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mini-clis/shared"
	"github.com/samber/lo"
)

// TODO: Think about the states of the app we have idle, editing, help, submit and submitting.
// There needs to be more inputs added to make up a form.
type model struct {
	terminalSizeManager shared.TerminalSizeManager
	inputModel          InputModel
	helpModel           help.Model
	submitted           bool
}

func InitialModel() model {

	return model{
		inputModel: TextInput("Name", "What is your name?").Focus(),
	}

}

type KeyMap struct {
	Quit  key.Binding
	Help  key.Binding
	Enter key.Binding
}

func NewKeyMap() KeyMap {

	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q ctrl+c", "Quit App"),
		),
		Enter: key.NewBinding(
			key.WithKeys(keys.Enter.String(), keys.Space.String()),
			key.WithHelp("enter or spacebar", "submit "),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
	}

}

func (km KeyMap) ShortHelp() []key.Binding {

	return []key.Binding{
		km.Enter,
		km.Quit,
	}

}

func (km KeyMap) FullHelp() [][]key.Binding {

	return [][]key.Binding{
		{km.Enter, km.Quit},
		{km.Help},
	}

}

// AfterSubmittedMsg is a message that is sent after a form has been submitted
// and helps control the flow of the form submission process
type AfterSumbitedMsg struct{}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	// TODO: Implement Help functionality here
	switch msg := message.(type) {

	case tea.KeyMsg:

		keyMap := NewKeyMap()

		switch {

		case key.Matches(msg, keyMap.Quit):

			return m, tea.Quit

		case key.Matches(msg, keyMap.Help):

			m.helpModel.ShowAll = !m.helpModel.ShowAll

		case key.Matches(msg, keyMap.Enter):

			m.submitted = true

			return m, func() tea.Msg {

				time.Sleep(time.Millisecond * 250)

				return AfterSumbitedMsg{}

			}

		}

	case tea.WindowSizeMsg:

		m.terminalSizeManager.SetTerminalSize(msg.Width, msg.Height)

	case AfterSumbitedMsg:

		m.submitted = false

		m.inputModel.Clear()

	}

	return m, m.inputModel.Update(message)

}

func (m model) View() string {

	AbsoluteCenter := func(terminalString string) string {

		terminalSize := m.terminalSizeManager.GetTerminalSize()

		return lipgloss.Place(
			terminalSize.Width(),
			terminalSize.Height(),
			lipgloss.Center,
			lipgloss.Center,
			terminalString,
		)
	}

	return AbsoluteCenter(
		lipgloss.JoinVertical(
			lipgloss.Center,
			shared.NewRows(2).
				Render(
					lo.If(
						m.submitted,
						fmt.Sprintf("Congradulations %s", m.inputModel.GetValue()),
					).
						Else(m.inputModel.View()),
					m.helpModel.View(NewKeyMap()),
				),
		),
	)
}

func (m model) Init() tea.Cmd {

	return nil

}
