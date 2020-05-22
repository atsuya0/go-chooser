package chooser

type keyBindingBufferCmd struct {
	function    func(*buffer)
	key         string
	description string
}

type keyBindingRenderCmd struct {
	function    func(*render, bool)
	key         string
	description string
}

var keyBindingBufferCmds = map[int]keyBindingBufferCmd{
	delete: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.deleteCharOnCursor()
		},
		key:         "delete",
		description: "Delete a character under cursor.",
	},
	controlD: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.deleteCharOnCursor()
		},
		key:         "controlD",
		description: "Delete a character under cursor.",
	},
	backspace: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.backwardDeleteChar()
		},
		key:         "backspace",
		description: "Delete a character before cursor.",
	},
	controlH: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.backwardDeleteChar()
		},
		key:         "controlH",
		description: "Delete a character before cursor.",
	},
	controlF: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.forwardChar()
		},
		key:         "controlF",
		description: "Move forward a character.",
	},
	controlB: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.backwardChar()
		},
		key:         "controlB",
		description: "Move backward a character.",
	},
	controlA: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.beginningOfLine()
		},
		key:         "controlA",
		description: "Go to the beginning of the line.",
	},
	controlE: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.endOfLine()
		},
		key:         "controlE",
		description: "Go to the end of the line.",
	},
	controlU: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.backwardKillLine()
		},
		key:         "controlU",
		description: "Kill characters from cursor current position to the beginning of the line.",
	},
	controlK: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.killLine()
		},
		key:         "controlK",
		description: "Kill characters from cursor current position to the end of the line.",
	},
	controlW: keyBindingBufferCmd{
		function: func(b *buffer) {
			b.backwardKillWord()
		},
		key:         "controlW",
		description: "Delete before a word.",
	},
}

var keyBindingRenderCmds = map[int]keyBindingRenderCmd{
	controlN: keyBindingRenderCmd{
		function: func(r *render, _ bool) {
			r.next()
		},
		key:         "controlN",
		description: "Move the cursor to the next line.",
	},
	controlP: keyBindingRenderCmd{
		function: func(r *render, _ bool) {
			r.previous()
		},
		key:         "controlP",
		description: "Move the cursor to the previous line.",
	},
	tab: keyBindingRenderCmd{
		function: func(r *render, isMultiple bool) {
			if isMultiple {
				r.updateRegister()
			}
		},
		key:         "tab",
		description: "Store the line on the cursor.",
	},
}
