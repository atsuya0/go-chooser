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
