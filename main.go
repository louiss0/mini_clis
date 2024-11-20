package main

import (
	"fmt"
	"log"
	"time"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

type TermialSize struct {
	width  int
	height int
}

func (ts TermialSize) Width() int {

	return ts.width

}

func (ts TermialSize) Height() int {

	return ts.height

}

type model struct {
	terminalSize TermialSize
	inputModel   InputModel
	Submitted    bool
}

func (m *model) SetTerminalSize(width, height int) {

	m.terminalSize = TermialSize{width, height}

}

func (m model) GetTerminalSize() TermialSize {

	return m.terminalSize
}

func initialModel() model {

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
			key.WithKeys(keys.Enter.String(), "ctrl+c"),
			key.WithHelp("q ctrl+c", ""),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
	}

}

type InputModel struct {
	textInput    textinput.Model
	label        string
	defaultValue string
}

func (im *InputModel) Clear() {

	im.textInput.SetValue("")

}

func (im *InputModel) GetValue() string {

	return im.textInput.Value()
}

func (im *InputModel) Reset() {

	if im.defaultValue != "" {

		im.textInput.SetValue(im.defaultValue)

		return

	}

	im.textInput.Reset()

}

func (im InputModel) Focus() InputModel {

	im.textInput.Focus()

	return im

}

func (im *InputModel) Update(msg tea.Msg) tea.Cmd {

	textInput, cmd := im.textInput.Update(msg)

	im.textInput = textInput

	return cmd
}

func (im *InputModel) SetCharacterLimit(characterLimit int) *InputModel {

	im.textInput.CharLimit = characterLimit

	return im

}

func (im *InputModel) SetDefaultValue(value string) *InputModel {

	if im.defaultValue != "" {

		return im

	}

	im.textInput.SetValue(value)

	im.defaultValue = value

	return im

}

func (im *InputModel) SetWidth(width int) *InputModel {

	im.textInput.Width = width

	return im
}

func (im InputModel) View() string {

	return fmt.Sprintf("%s\n%s", im.label, im.textInput.View())

}

func TextInput(label, placeholder string) InputModel {

	textInput := textinput.New()

	textInput.CharLimit = 25
	textInput.TextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: "#1e293b",
			Dark:  "#f1f5f9",
		})

	textInput.Cursor.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22d3ee"))

	textInput.Placeholder = placeholder

	return InputModel{textInput, label, ""}

}

type AfterSumbitedMsg struct{}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := message.(type) {

	case tea.KeyMsg:

		keyMap := NewKeyMap()

		switch {

		case key.Matches(msg, keyMap.Quit):

			return m, tea.Quit

		case key.Matches(msg, keyMap.Help):

		case key.Matches(msg, keyMap.Enter):

			m.Submitted = true

			return m, func() tea.Cmd {

				return func() tea.Msg {

					time.Sleep(time.Millisecond * 250)

					return AfterSumbitedMsg{}

				}

			}()

		}

	case tea.WindowSizeMsg:

		m.SetTerminalSize(msg.Width, msg.Height)

	case AfterSumbitedMsg:

		m.Submitted = false

		m.inputModel.Clear()

	}

	cmd := m.inputModel.Update(message)

	return m, cmd

}

func (m model) View() string {

	AbsoluteCenter := func(terminalString string) string {

		terminalSize := m.GetTerminalSize()

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
			lo.If(
				m.Submitted,
				fmt.Sprintf("Congradulations %s", m.inputModel.GetValue()),
			).Else(m.inputModel.View()),
		),
	)
}

func (m model) Init() tea.Cmd {

	return nil

}

func main() {

	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {

		log.Fatal(err)

	}

}
