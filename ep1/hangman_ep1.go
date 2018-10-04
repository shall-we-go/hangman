package ep1

import (
	"strings"
)

// Input: hangman('bigbears', ['a','e','i','o','u','c','d','p','r','k','l','j','h'])
// Expected output: false
// Input: hangman('bigbears', ['s', 'r', 'a', 'e', 'b', 'g', 'i', 'b'])
// Expected output: true
func Hangman(secretWord string, letters []string) bool {
	for _, letter := range letters {
		secretWord = strings.Replace(secretWord, letter, "", -1)
		if secretWord == "" {
			return true
		}
	}
	return false
}
