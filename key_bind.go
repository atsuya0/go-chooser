package chooser

type bufferKeyBindingCmd struct {
	function    func(*buffer)
	key         string
	description string
}

type renderKeyBindingCmd struct {
	function    func(*render, bool)
	key         string
	description string
}

type keyBindingCmd struct {
	key         string
	description string
}

var bufferKeyBindingCmds = map[int]bufferKeyBindingCmd{
	delete: {
		function: func(b *buffer) {
			b.deleteCharOnCursor()
		},
		key:         "delete",
		description: "Delete a character under cursor.",
	},
	controlD: {
		function: func(b *buffer) {
			b.deleteCharOnCursor()
		},
		key:         "controlD",
		description: "Delete a character under cursor.",
	},
	backspace: {
		function: func(b *buffer) {
			b.backwardDeleteChar()
		},
		key:         "backspace",
		description: "Delete a character before cursor.",
	},
	controlH: {
		function: func(b *buffer) {
			b.backwardDeleteChar()
		},
		key:         "controlH",
		description: "Delete a character before cursor.",
	},
	controlF: {
		function: func(b *buffer) {
			b.forwardChar()
		},
		key:         "controlF",
		description: "Move forward a character.",
	},
	controlB: {
		function: func(b *buffer) {
			b.backwardChar()
		},
		key:         "controlB",
		description: "Move backward a character.",
	},
	controlA: {
		function: func(b *buffer) {
			b.beginningOfLine()
		},
		key:         "controlA",
		description: "Go to the beginning of the line.",
	},
	controlE: {
		function: func(b *buffer) {
			b.endOfLine()
		},
		key:         "controlE",
		description: "Go to the end of the line.",
	},
	controlU: {
		function: func(b *buffer) {
			b.backwardKillLine()
		},
		key:         "controlU",
		description: "Kill characters from cursor current position to the beginning of the line.",
	},
	controlK: {
		function: func(b *buffer) {
			b.killLine()
		},
		key:         "controlK",
		description: "Kill characters from cursor current position to the end of the line.",
	},
	controlW: {
		function: func(b *buffer) {
			b.backwardKillWord()
		},
		key:         "controlW",
		description: "Delete before a word.",
	},
}

var renderKeyBindingCmds = map[int]renderKeyBindingCmd{
	controlN: {
		function: func(r *render, _ bool) {
			r.next()
		},
		key:         "controlN",
		description: "Move the cursor to the next line.",
	},
	controlP: {
		function: func(r *render, _ bool) {
			r.previous()
		},
		key:         "controlP",
		description: "Move the cursor to the previous line.",
	},
	tab: {
		function: func(r *render, isMultiple bool) {
			if isMultiple {
				r.holdCompletion()
			}
		},
		key:         "tab",
		description: "Store the line on the cursor.",
	},
}

var keyBindingCmds = map[int]keyBindingCmd{
	enter: {
		key:         "enter",
		description: "Choose the line on the cursor. Or choose the stored lines.",
	},
	controlC: {
		key:         "controlC",
		description: "Cancel.",
	},
}
