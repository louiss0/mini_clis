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

func NewTerminalSize(width, height int) TermialSize {

	return TermialSize{width, height}
}
