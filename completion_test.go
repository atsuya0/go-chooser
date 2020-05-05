package choice

import "testing"

func TestCompletionNew(t *testing.T) {
	patterns := []*completion{
		{
			suggestions: []string{"test1", "test2", "test3"},
			indexes:     []int{2, 6, 8},
			target:      0,
		},
		{
			suggestions: []string{},
			indexes:     []int{},
			target:      -1,
		},
	}

	for _, pattern := range patterns {
		completion := newCompletion(pattern.suggestions, pattern.indexes)
		if completion.target != pattern.target {
			t.Errorf("output: %v, request: %v", completion, pattern)
		}
	}
}

func TestCompletionNext(t *testing.T) {
	patterns := []struct {
		completion    *completion
		requestTarget int
	}{
		{
			completion: &completion{
				suggestions: []string{"test1", "test2", "test3"}, target: 0,
			},
			requestTarget: 1,
		},
		{
			completion: &completion{
				suggestions: []string{}, target: -1,
			},
			requestTarget: -1,
		},
	}

	for _, pattern := range patterns {
		pattern.completion.next()
		if pattern.completion.target != pattern.requestTarget {
			t.Errorf("output: %d, request: %d",
				pattern.completion.target, pattern.requestTarget)
		}
	}
}

func TestCompletionPrevious(t *testing.T) {
	patterns := []struct {
		completion    *completion
		requestTarget int
	}{
		{
			completion: &completion{
				suggestions: []string{"test1", "test2", "test3"}, target: 2,
			},
			requestTarget: 1,
		},
		{
			completion: &completion{
				suggestions: []string{}, target: -1,
			},
			requestTarget: -1,
		},
	}

	for _, pattern := range patterns {
		pattern.completion.previous()
		if pattern.completion.target != pattern.requestTarget {
			t.Errorf("output: %d, request: %d",
				pattern.completion.target, pattern.requestTarget)
		}
	}
}
