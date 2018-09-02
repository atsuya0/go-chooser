package selector

var keyBindCmds = map[int]func(*buffer){
	delete: func(b *buffer) {
		b.deleteChar()
	},
	controlD: func(b *buffer) {
		b.deleteChar()
	},
	backspace: func(b *buffer) {
		b.backwardDeleteChar()
	},
	controlH: func(b *buffer) {
		b.backwardDeleteChar()
	},
	controlF: func(b *buffer) {
		b.forwardChar()
	},
	controlB: func(b *buffer) {
		b.backwardChar()
	},
	controlA: func(b *buffer) {
		b.beginningOfLine()
	},
	controlE: func(b *buffer) {
		b.endOfLine()
	},
	controlU: func(b *buffer) {
		b.backwardKillLine()
	},
	controlK: func(b *buffer) {
		b.killLine()
	},
	controlW: func(b *buffer) {
		b.backwardKillWord()
	},
}
