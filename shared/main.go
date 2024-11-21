package shared

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

var NewTermialSizeManager = new(TermialSizeManager)

type TermialSizeManager struct {
	terminalSize TermialSize
}

func (tsm *TermialSizeManager) SetTerminalSize(width, height int) {

	tsm.terminalSize = TermialSize{width, height}

}

func (tsm TermialSizeManager) GetTerminalSize() TermialSize {

	return tsm.terminalSize
}
