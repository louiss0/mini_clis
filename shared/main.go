package shared

import (
	"fmt"
	"maps"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

type BackedEnum[
	ValueType interface{ string | int },
	Map ~map[string]ValueType,
] struct {
	structure Map
}

func (e BackedEnum[ValueType, Map]) Structure() Map {

	return maps.Clone(e.structure)

}

func (e BackedEnum[ValueType, Map]) Validate(input ValueType) bool {

	for _, value := range e.Values() {

		if value == input {

			return true
		}

	}

	return false

}

func (e BackedEnum[ValueType, Map]) Parse(input ValueType) error {

	for _, value := range e.Values() {

		if value == input {

			return nil
		}

	}

	return fmt.Errorf("invalid enum value %v", input)

}

func (e BackedEnum[ValueType, Map]) Values() []ValueType {

	slice := []ValueType{}

	structValues := maps.Values(e.structure)

	for value := range structValues {

		slice = append(slice, value)

	}

	return slice

}

type TerminalSize struct {
	width  int
	height int
}

func (ts TerminalSize) Width() int {

	return ts.width

}

func (ts TerminalSize) Height() int {

	return ts.height

}

var NewTerminalSizeManager = new(TerminalSizeManager)

type TerminalSizeManager struct {
	terminalSize TerminalSize
}

func (tsm *TerminalSizeManager) SetTerminalSize(width, height int) {

	tsm.terminalSize = TerminalSize{width, height}

}

func (tsm TerminalSizeManager) GetTerminalSize() TerminalSize {

	return tsm.terminalSize
}

var gapRenderer = lipgloss.NewStyle()

type GapNumber = interface{ int | float64 }

func HorizontalGap[T GapNumber](size T) string {

	return gapRenderer.UnsetHeight().Width(int(size)).Render()

}

func VerticalGap[T GapNumber](size T) string {

	return gapRenderer.UnsetWidth().Height(int(size)).Render()

}

type Rows[T GapNumber] struct {
	gap T
}

func NewRows[T GapNumber](gap T) Rows[T] {

	return Rows[T]{gap}
}

func (r Rows[T]) Render(renderStrings ...string) string {

	contents := lo.Map(renderStrings, func(item string, index int) string {

		return item
	})

	return strings.Join(contents, VerticalGap(r.gap))
}

type Columns[T GapNumber] struct {
	gap T
}

func NewColumns[T GapNumber](gap T) Columns[T] {

	return Columns[T]{gap}
}

func (c Columns[T]) Render(gap T, renderStrings ...string) string {

	contents := lo.Map(renderStrings, func(item string, index int) string {

		return item
	})

	return strings.Join(contents, HorizontalGap(gap))

}
