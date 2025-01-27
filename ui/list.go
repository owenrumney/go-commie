package ui

import (
	"fmt"
)

func ChooseFromList(prompt string, items []string) (int, string, error) {

	fmt.Printf("\n%s: \n\n", prompt)

	for _, item := range items {
		fmt.Println(item)
	}

	moveCursorUp(len(items))

	currentPos := 0

keyInput:
	for {
		keyCode, err := getKeyInput()
		if err != nil {
			return 0, "", err
		}
		switch keyCode {
		case UP:
			if currentPos > 0 {
				currentPos--
				moveCursorUp(1)
			}
		case DOWN:
			if currentPos < len(items)-1 {
				currentPos++
				moveCursorDown(1)
			}
		case RETURN:
			break keyInput
		case ESCAPE:
			return 0, "", fmt.Errorf("user cancelled")
		}

	}

	resetPrompt(len(items) - currentPos)
	return currentPos, items[currentPos], nil
}

func resetPrompt(rowPosition int) {
	moveCursorDown(rowPosition + ROW_OFFSET - 1)
	clearLine()
	moveCursorToColumn(DEFAULT_COLUMN)
}
