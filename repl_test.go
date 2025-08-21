package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HARRY POTTEr pS1",
			expected: []string{"harry", "potter", "ps1"},
		},
		{
			input:    " micHaEl  JAckSoN",
			expected: []string{"michael", "jackson"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("expected length: %d, got: %d", len(c.expected), len(actual))
			return
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected word: %s, got: %s", expectedWord, word)
				return
			}
		}
	}
}
