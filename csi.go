package chooser

import (
	"fmt"
)

func cursorUp(num int) string {
	return fmt.Sprintf("\x1b[%dA", num)
}

func cursorForward(num int) string {
	return fmt.Sprintf("\x1b[%dG", num)
}

// Clear from the cursor position to the end.
func clearCursorEnd() string {
	return fmt.Sprint("\x1b[J")
}

// Clear the screen.
func clearScreen() string {
	return cursorForward(1) + clearCursorEnd()
}

func restoreCursorPosition(row, col int) string {
	return cursorUp(row) + cursorForward(col)
}
