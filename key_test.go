package choice

import (
	"testing"
)

func TestGetKey(t *testing.T) {
	keys := []*key{
		{id: controlA, code: []byte{0x1}},
		{id: delete, code: []byte{0x1b, 0x5b, 0x33, 0x7e}},
		{id: displayable, code: []byte{0x20}}, // space
		{id: displayable, code: []byte{0x7e}}, // ~
		{id: ignore, code: []byte{0x80}},
	}

	for _, key := range keys {
		if id := getKey(key.code); id != key.id {
			t.Errorf("[Key] output: %d, request: %d:%s", id, key.id, key.code)
		}
	}
}
