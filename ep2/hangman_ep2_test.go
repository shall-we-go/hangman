package ep2

import (
	"reflect"
	"testing"
)

func TestReactiveHangman(t *testing.T) {
	cases := []struct {
		inSecretWord string
		inLetters    []string
		want         Result
	}{
		{
			inSecretWord: "bigbear",
			inLetters:    []string{},
			want: Result{
				status:           InProgress,
				selectedLetters:  []string{},
				lifeLeft:         7,
				secretWordLength: 7,
				knownSecretWord:  "_______",
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b"},
			want: Result{
				status:           InProgress,
				selectedLetters:  []string{"b"},
				lifeLeft:         7,
				secretWordLength: 7,
				knownSecretWord:  "b__b___",
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o"},
			want: Result{
				status:           InProgress,
				selectedLetters:  []string{"b", "o"},
				lifeLeft:         6,
				secretWordLength: 7,
				knownSecretWord:  "b__b___",
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o", "i", "g", "a", "e", "y", "r"},
			want: Result{
				status:           Won,
				selectedLetters:  []string{"b", "o", "i", "g", "a", "e", "y", "r"},
				lifeLeft:         5,
				secretWordLength: 7,
				knownSecretWord:  "bigbear",
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o", "a", "e", "n", "u", "t", "z", "x", "v"},
			want: Result{
				status:           Lose,
				selectedLetters:  []string{"b", "o", "a", "e", "n", "u", "t", "z", "x", "v"},
				lifeLeft:         0,
				secretWordLength: 7,
				knownSecretWord:  "b__bea_",
			},
		},
	}
	for _, c := range cases {
		letterCh := make(chan string)
		gotCh := ReactiveHangman(c.inSecretWord, letterCh)
		go func() {
			for _, letter := range c.inLetters {
				letterCh <- letter
			}
			//Close once no more new letter
			close(letterCh)
		}()
		var got Result
		for got = range gotCh {

		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Hangman(%q, %q) == %#v, want %#v", c.inSecretWord, c.inLetters, got, c.want)
		}
	}
}
