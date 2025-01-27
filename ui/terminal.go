package ui

import (
	"fmt"
)

func moveCursorToColumn(column int) {
	fmt.Printf("\033[%dG", column+1)
}

func moveCursorDown(rows int) {
	fmt.Printf("\033[%dB", rows)
}

func moveCursorUp(rows int) {
	if rows > 0 {
		fmt.Printf("\033[%dA", rows)
	}
}

func clearLine() {
	fmt.Printf("\033[2K\r")
}
