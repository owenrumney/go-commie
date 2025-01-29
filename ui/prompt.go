package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetInput(question string) string {
	clearLine()
	fmt.Printf("\n%s: ", question)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil || len(input) <= 1 {
		return ""
	}
	s := input[:len(input)-1]
	if s[len(s)-1] == '\r' {
		s = input[:len(input)-1]
	}
	return s
}

func GetMultilineInput(question string) string {
	clearLine()
	fmt.Printf("\n%s (finish with two empty lines):\n\n", question)

	reader := bufio.NewReader(os.Stdin)
	var lines []string
	consecutiveEmptyLines := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			consecutiveEmptyLines++
		} else {
			consecutiveEmptyLines = 0
		}
		if consecutiveEmptyLines == 2 {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func YesNoQuestion(question string, defaultToYes bool) (bool, error) {
	if defaultToYes {
		question += " [Y/n]: "
	} else {
		question += " [y/N]: "
	}

	clearLine()
	fmt.Printf("%s", question)

	for {
		keyCode, err := getKeyInput()
		if err != nil {
			return false, nil
		}

		switch keyCode {
		case Y:
			return true, nil
		case N:
			return false, nil
		case RETURN:
			return defaultToYes, nil
		case ESCAPE:
			return false, fmt.Errorf("user cancelled")
		}
	}
}
