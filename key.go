package selector

import "bytes"

type key struct {
	key  int
	code []byte
}

const (
	escape = iota
	controlA
	controlB
	controlC
	controlD
	controlE
	controlF
	controlG
	controlH
	controlI
	controlJ
	controlK
	controlL
	controlM
	controlN
	controlO
	controlP
	controlQ
	controlR
	controlS
	controlT
	controlU
	controlV
	controlW
	controlX
	controlY
	controlZ
	enter
	tab
	backTab
	delete
	backspace

	ignore
	displayable
)

var keys = []*key{
	{key: controlA, code: []byte{0x1}},
	{key: controlB, code: []byte{0x2}},
	{key: controlC, code: []byte{0x3}},
	{key: controlD, code: []byte{0x4}},
	{key: controlE, code: []byte{0x5}},
	{key: controlF, code: []byte{0x6}},
	{key: controlG, code: []byte{0x7}},
	{key: controlH, code: []byte{0x8}},
	// tab: {key: controlI, code: []byte{0x9}},
	// enter: {key: controlJ, code: []byte{0xa}},
	{key: controlK, code: []byte{0xb}},
	{key: controlL, code: []byte{0xc}},
	{key: controlM, code: []byte{0xd}},
	{key: controlN, code: []byte{0xe}},
	{key: controlO, code: []byte{0xf}},
	{key: controlP, code: []byte{0x10}},
	{key: controlQ, code: []byte{0x11}},
	{key: controlR, code: []byte{0x12}},
	{key: controlS, code: []byte{0x13}},
	{key: controlT, code: []byte{0x14}},
	{key: controlU, code: []byte{0x15}},
	{key: controlV, code: []byte{0x16}},
	{key: controlW, code: []byte{0x17}},
	{key: controlX, code: []byte{0x18}},
	{key: controlY, code: []byte{0x19}},
	{key: controlZ, code: []byte{0x1a}},
	{key: escape, code: []byte{0x1b}},
	{key: enter, code: []byte{0xa}},
	{key: backspace, code: []byte{0x7f}},
	{key: delete, code: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{key: tab, code: []byte{0x9}},
	{key: backTab, code: []byte{0x1b, 0x5b, 0x5a}},
}

func getKey(b []byte) int {
	// Predefined key
	for _, key := range keys {
		if bytes.Equal(key.code, b) {
			return key.key
		}
	}
	// Displayable characters.
	if len(b) == 1 && 0x21 <= b[0] && b[0] <= 0x7e {
		return displayable
	}

	return ignore
}
