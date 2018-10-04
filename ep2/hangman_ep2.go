package ep2

import (
	"strings"
)

type Status int

const (
	InProgress Status = iota
	Won
	Lose
)

const MaxLife = 7
const BlindLetter = "_"

type Result struct {
	status           Status
	selectedLetters  []string
	lifeLeft         int
	secretWordLength int
	knownSecretWord  string
}

func ReactiveHangman(secretWord string, letterCh chan string) chan Result {
	resultCh := make(chan Result)
	go func() {
		//Initial game
		result := Result{
			status:           InProgress,
			selectedLetters:  []string{},
			lifeLeft:         MaxLife,
			secretWordLength: len(secretWord),
			knownSecretWord:  strings.Repeat(BlindLetter, len(secretWord)),
		}
		resultCh <- result

		//Received a new letter?
		for letter := range letterCh {
			//Recorded the letter
			result.selectedLetters = append(result.selectedLetters, letter)
			//Found new letter?
			found := false
			for i, secretLetter := range secretWord {
				if secretLetter == []rune(letter)[0] {
					found = true
					result.knownSecretWord = result.knownSecretWord[:i] + letter + result.knownSecretWord[i+1:]
				}
			}
			// Not Found?
			if !found {
				result.lifeLeft--
			}
			//Won?
			if strings.Index(result.knownSecretWord, BlindLetter) == -1 {
				result.status = Won
			}
			//Lost?
			if result.lifeLeft == 0 {
				result.status = Lose
			}
			resultCh <- result
		}
		//Close once no more new letter
		close(resultCh)
	}()
	return resultCh
}
