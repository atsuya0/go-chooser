package choice

import (
	"fmt"
)

func cursorUp(num int) string {
	return fmt.Sprintf("\x1b[%dA", num)
}

// Set the column at the cursor position.
func setColCursor(num int) string {
	return fmt.Sprintf("\x1b[%dG", num)
}

// Clear from the cursor position to the end.
func clearCursorEnd() string {
	return "\x1b[J"
}

// Clear the screen.
func clear() {
	fmt.Print(setColCursor(1), clearCursorEnd())
}
