package ep3

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	InProgress Status = "in-progress"
	Won        Status = "won"
	Lose       Status = "lose"
)

type Event string

const (
	TimeSpent Event = "time-spent"
	GameOver  Event = "game-over"
)

const MaxLife = 7
const BlindLetter = "_"
const Timeout = 5
const NumLettersPerHint = 5

type Result struct {
	Id               string
	Status           Status
	SelectedLetters  []string
	LifeLeft         int
	SecretWordLength int
	KnownSecretWord  string
}

type Timer struct {
	Id    string
	Event Event
	Data  TimerData
}

type TimerData struct {
	TimeLeft int
	LifeLeft int
}

type Game struct {
	Id              string
	Result          Result
	Timer           Timer
	letterCh        chan string
	resultByGuessCh chan Result
	timerChs        []chan Timer
	SecretWord      string
}

var Games = make(map[string]*Game)

type RandomWordService interface {
	RandomWord() string
}

type RandomHintLettersService interface {
	RandomHintLetters(secretWord string) []string
}

type TimerIdService interface {
	GetTimerId() string
}

func (this *Game) create(randomWordService RandomWordService, randomHintLettersService RandomHintLettersService, timerIdService TimerIdService) chan Result {
	createdResult := make(chan Result)
	letterCh := make(chan string)
	this.letterCh = letterCh
	resultByGuessCh := make(chan Result)
	this.resultByGuessCh = resultByGuessCh
	timerCh := make(chan Timer)
	this.SecretWord = randomWordService.RandomWord()
	log.Println(this.SecretWord)
	hintLetters := randomHintLettersService.RandomHintLetters(this.SecretWord)
	resultCh := this.ReactiveHangman(this.SecretWord, hintLetters, letterCh, timerCh, resultByGuessCh)
	awaitingFirstResult := true
	go func() {
	GameLoop:
		for {
			select {
			case result := <-resultCh:
				log.Println(result)
				if awaitingFirstResult {
					awaitingFirstResult = false
					createdResult <- result
				}
				this.Result = result
				if result.Status == Won || result.Status == Lose {
					break GameLoop
				}
			case timer := <-timerCh:
				log.Println(timer)
				this.Timer = timer
				timer.Id = timerIdService.GetTimerId()
				for _, clientTimerCh := range this.timerChs {
					clientTimerCh <- timer
				}
			}
		}
	}()
	return createdResult
}

func (this *Game) ReactiveHangman(secretWord string, hintLetters []string, letterCh chan string, timerCh chan Timer, resultByGuessCh chan Result) chan Result {
	resultCh := make(chan Result)
	go func() {
		//Initial game
		result := Result{
			Id:               this.Id,
			Status:           InProgress,
			SelectedLetters:  []string{},
			LifeLeft:         MaxLife,
			SecretWordLength: len(secretWord),
			KnownSecretWord:  strings.Repeat(BlindLetter, len(secretWord)),
		}
		timer := Timer{
			Event: TimeSpent,
			Data: TimerData{
				TimeLeft: Timeout,
				LifeLeft: result.LifeLeft,
			},
		}

		ticker := time.NewTicker(time.Second)
		mapLetters := map[string]bool{}

		kill := func() {
			result.LifeLeft--
			timer.Data.LifeLeft = result.LifeLeft
		}

		lose := func() {
			result.Status = Lose
			timer.Event = GameOver
			ticker.Stop()
		}

		updateResultByLetter := func(letter string) bool {
			//In previous letters ?
			if mapLetters[letter] {
				return false
			}
			mapLetters[letter] = true
			timer.Data.TimeLeft = Timeout
			//Recorded the letter
			result.SelectedLetters = append(result.SelectedLetters, letter)
			//Found new letter?
			found := false
			for i, secretLetter := range secretWord {
				if secretLetter == []rune(letter)[0] {
					found = true
					result.KnownSecretWord = result.KnownSecretWord[:i] + letter + result.KnownSecretWord[i+1:]
				}
			}
			// Not Found?
			if !found {
				kill()
			}
			if result.LifeLeft == 0 {
				//Lost?
				lose()
			} else if strings.Index(result.KnownSecretWord, BlindLetter) == -1 {
				//Won?
				result.Status = Won
				timer.Event = GameOver
				ticker.Stop()
			}
			return true
		}

		//Update result by hint letters
		for _, hintLetter := range hintLetters {
			updateResultByLetter(hintLetter)
		}

		resultCh <- result
		timerCh <- timer

		for {
			select {
			//Received a new letter?
			case letter := <-letterCh:
				if updateResultByLetter(letter) {
					timerCh <- timer
				}
				resultCh <- result
				go func() {
					resultByGuessCh <- result
				}()
			//Received counter?
			case <-ticker.C:
				timer.Data.TimeLeft--
				if timer.Data.TimeLeft == 0 {
					kill()
					//Lost?
					if result.LifeLeft == 0 {
						lose()
						//Send TimeLeft=0 Event=GameOver
						timerCh <- timer
					} else {
						//Send TimeLeft=0
						timerCh <- timer
						//Send TimeLeft=5
						timer.Data.TimeLeft = Timeout
						time.Sleep(time.Millisecond)
						timerCh <- timer
					}
					resultCh <- result
				} else {
					timerCh <- timer
				}
			}
		}
	}()
	return resultCh
}

func (this *Game) RandomWord() string {
	rand.Seed(time.Now().UnixNano())
	words := []string{
		"adventurous",
		"courageous",
		"extramundane",
		"generous",
		"intransigent",
		"sympathetic",
		"vagarious",
		"witty",
	}
	return words[rand.Intn(len(words))]
}

func (this *Game) GetTimerId() string {
	return strconv.Itoa(int(time.Now().UTC().Unix()))
}

func (this *Game) RandomHintLetters(secretWord string) []string {
	numHint := int(math.Ceil(float64(len(secretWord)) / float64(NumLettersPerHint)))
	hintLetters := make([]string, numHint)
	lettersMap := map[string]bool{}
	letters := []string{}
	for _, r := range secretWord {
		c := string(r)
		if !lettersMap[c] {
			lettersMap[c] = true
			letters = append(letters, c)
		}
	}
	for i := 0; i < numHint; i++ {
		rand.Seed(time.Now().UnixNano())
		rIndex := rand.Intn(len(letters))
		hintLetters[i] = letters[rIndex]
		letters = append(letters[:rIndex], letters[rIndex+1:]...)
	}
	return hintLetters
}

func CreateNewGame() Result {
	game := Game{
		Id: "gameId-" + strconv.Itoa(len(Games)+1),
	}
	Games[game.Id] = &game
	resultCh := game.create(&game, &game, &game)
	return <-resultCh
}

func GetGameInfo(id string) Result {
	game := Games[id]
	if game == nil {
		return Result{}
	}
	return game.Result
}

func GuessLetter(id string, letter string) Result {
	game := Games[id]
	if game == nil {
		return Result{}
	} else if game.Result.Status != InProgress {
		return game.Result
	}
	game.letterCh <- letter
	return <-game.resultByGuessCh
}

func (this *Game) SubscibeTimer() chan Timer {
	timerCh := make(chan Timer)
	this.timerChs = append(this.timerChs, timerCh)
	return timerCh
}
