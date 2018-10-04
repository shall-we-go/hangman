package ep1

import (
	"testing"
)

func TestHangman(t *testing.T) {
	cases := []struct {
		inSecretWord string
		inLetters    []string
		want         bool
	}{
		{
			inSecretWord: "abc",
			inLetters:    []string{"a", "b"},
			want:         false,
		},
		{
			inSecretWord: "abc",
			inLetters:    []string{"a", "b", "c"},
			want:         true,
		},
		{
			inSecretWord: "bigbears",
			inLetters:    []string{"a", "e", "i", "o", "u", "c", "d", "p", "r", "k", "l", "j", "h"},
			want:         false,
		},
		{
			inSecretWord: "bigbears",
			inLetters:    []string{"s", "r", "a", "e", "b", "g", "i", "b"},
			want:         true,
		},
	}
	for _, c := range cases {
		got := Hangman(c.inSecretWord, c.inLetters)
		if got != c.want {
			t.Errorf("Hangman(%q, %q) == %v, want %v", c.inSecretWord, c.inLetters, got, c.want)
		}
	}
}
