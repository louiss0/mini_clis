package counter

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
)

const INCREMENT = "increment"

const DECREMENT = "decrement"

type KeyMap struct {
	Up   key.Binding
	Down key.Binding
	Help key.Binding
	Quit key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "increase count"),      // actual keybindings
		key.WithHelp("↑/k", "Increment Counter"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "decrease count"),
		key.WithHelp("↓/j", "Decrement Counter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},   // first column
		{k.Help, k.Quit}, // second column
	}
}

type model struct {
	count  int
	width  int
	height int
}

func (m model) Init() tea.Cmd {

	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:

		switch {

		case key.Matches(msg, DefaultKeyMap.Up):

			m.count++

		case key.Matches(msg, DefaultKeyMap.Down):
			if m.count <= 0 {

				m.count = 0

			} else {

				m.count--

			}
		case key.Matches(msg, DefaultKeyMap.Quit):

			return m, tea.Quit

		}

	case tea.MouseMsg:

		if msg.Action == tea.MouseActionRelease || msg.Button == tea.MouseButtonLeft {

			if zone.Get(DECREMENT).InBounds(msg) {

				if m.count <= 0 {

					m.count = 0
				} else {

					m.count--
				}

			}

			if zone.Get(INCREMENT).InBounds(msg) {

				m.count++

			}

		}

		return m, nil

	}
	return m, nil

}

func HelpMenu() string {

	helpMenu := help.New()

	helpMenu.ShowAll = true

	return helpMenu.View(DefaultKeyMap)

}

func (m model) View() string {

	pterm.DefaultArea.WithCenter()

	title := pterm.DefaultBasicText.WithStyle(
		pterm.NewStyle(pterm.Bold, pterm.FgCyan),
	).Sprint("Counter App")

	counterBox := lipgloss.JoinHorizontal(
		lipgloss.Center,
		IncrementButton(),
		lipgloss.NewStyle().Width(2).Render(),
		pterm.FgGreen.Sprintf("Count: %d", m.count),
		lipgloss.NewStyle().Width(2).Render(),
		DecrementButton(),
	)

	return zone.Scan(
		lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				title,
				lipgloss.NewStyle().Width(2).Render(),
				counterBox,
				lipgloss.NewStyle().Width(1).Render(),
				HelpMenu(),
			),
		),
	)
}

func initialModel() tea.Model {

	return new(model)

}

func IncrementButton() string {

	return zone.Mark(
		INCREMENT,
		lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder(), true).
			Background(lipgloss.Color("92")).
			Foreground(lipgloss.Color("47")).
			Padding(1, 3).
			Render(lo.Capitalize(INCREMENT)),
	)
}
func DecrementButton() string {

	return zone.Mark(
		DECREMENT,
		lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder(), true).
			Background(lipgloss.Color("31")).
			Foreground(lipgloss.Color("47")).
			Padding(1, 3).
			Render(lo.Capitalize(DECREMENT)),
	)
}

func main() {

	zone.NewGlobal()

	program := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := program.Run(); err != nil {

		log.Fatal(err)

	}

}
