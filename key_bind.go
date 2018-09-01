package selector

var keyBindCmds = map[int]func(*buffer){
	delete: func(b *buffer) {
		b.deleteChar()
	},
	backspace: func(b *buffer) {
		b.deleteBeforeChar()
	},
	controlF: func(b *buffer) {
		b.goRightChar()
	},
	controlB: func(b *buffer) {
		b.goLeftChar()
	},
	controlA: func(b *buffer) {
		b.goLineBeginning()
	},
	controlE: func(b *buffer) {
		b.goLineEnd()
	},
	controlU: func(b *buffer) {
		b.backwardKillLine()
	},
}
