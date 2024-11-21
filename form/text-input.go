package form

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
