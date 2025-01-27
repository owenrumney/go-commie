package ui

import "github.com/pkg/term"

const (
	UP     = 65
	DOWN   = 66
	ESCAPE = 27
	RETURN = 13

	N = 110
	Y = 121

	ROW_OFFSET     = 2
	DEFAULT_COLUMN = 0
)

func getKeyInput() (keyCode int, err error) {
	t, err := term.Open("/dev/tty")
	if err != nil {
		return 0, err
	}
	err = term.RawMode(t)
	if err != nil {
		return 0, err
	}
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return 0, err
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		switch bytes[2] {
		case UP, DOWN:
			keyCode = int(bytes[2])
		}
	} else if numRead == 1 {
		switch bytes[0] {
		case ESCAPE, RETURN, N, Y:
			keyCode = int(bytes[0])
		}
	}
	t.Restore()
	t.Close()
	return
}
