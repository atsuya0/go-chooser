package choice

type completion struct {
	target      int
	suggestions []string
	indexes     []int
}

func (c *completion) length() int {
	return len(c.suggestions)
}

func (c *completion) next() {
	if c.target < 0 {
		return
	}

	if c.target+1 < c.length() {
		c.target += 1
	}
}

func (c *completion) previous() {
	if c.target < 0 {
		return
	}

	if 0 <= c.target-1 {
		c.target -= 1
	}
}

func (c *completion) getIndex() int {
	if c.target < 0 {
		return -1
	}
	return c.indexes[c.target]
}

func newCompletion(suggestions []string, indexes []int) *completion {
	var idx int
	if len(suggestions) == 0 {
		idx = -1
	} else {
		idx = 0
	}
	return &completion{
		target:      idx,
		suggestions: suggestions,
		indexes:     indexes,
	}
}
