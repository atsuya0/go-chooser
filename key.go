package chooser

import "bytes"

type key struct {
	id   int
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
	backspace
	delete

	displayable
	ignore
	question
)

var keys = []*key{
	{id: controlA, code: []byte{0x1}},
	{id: controlB, code: []byte{0x2}},
	{id: controlC, code: []byte{0x3}},
	{id: controlD, code: []byte{0x4}},
	{id: controlE, code: []byte{0x5}},
	{id: controlF, code: []byte{0x6}},
	{id: controlG, code: []byte{0x7}},
	{id: controlH, code: []byte{0x8}},
	{id: controlK, code: []byte{0xb}},
	{id: controlL, code: []byte{0xc}},
	{id: controlM, code: []byte{0xd}},
	{id: controlN, code: []byte{0xe}},
	{id: controlO, code: []byte{0xf}},
	{id: controlP, code: []byte{0x10}},
	{id: controlQ, code: []byte{0x11}},
	{id: controlR, code: []byte{0x12}},
	{id: controlS, code: []byte{0x13}},
	{id: controlT, code: []byte{0x14}},
	{id: controlU, code: []byte{0x15}},
	{id: controlV, code: []byte{0x16}},
	{id: controlW, code: []byte{0x17}},
	{id: controlX, code: []byte{0x18}},
	{id: controlY, code: []byte{0x19}},
	{id: controlZ, code: []byte{0x1a}},
	{id: escape, code: []byte{0x1b}},
	{id: enter, code: []byte{0xa}},
	{id: tab, code: []byte{0x9}},
	{id: backspace, code: []byte{0x7f}},
	{id: delete, code: []byte{0x1b, 0x5b, 0x33, 0x7e}},
	{id: question, code: []byte{0x3f}},
}

func getKey(b []byte) int {
	// Predefined key
	for _, key := range keys {
		if bytes.Equal(key.code, b) {
			return key.id
		}
	}
	// Displayable characters.(space and graphic character.)
	if len(b) == 1 && 0x20 <= b[0] && b[0] <= 0x7e {
		return displayable
	}

	return ignore
}
