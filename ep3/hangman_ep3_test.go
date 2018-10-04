package ep3

import (
	"log"
	"reflect"
	"testing"
)

func TestReactiveHangman(t *testing.T) {
	cases := []struct {
		inSecretWord string
		inLetters    []string
		result       Result
		timer        Timer
	}{
		{
			inSecretWord: "bigbear",
			inLetters:    []string{},
			result: Result{
				Id:               "gameId-0",
				Status:           Lose,
				SelectedLetters:  []string{},
				LifeLeft:         0,
				SecretWordLength: 7,
				KnownSecretWord:  "_______",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 0,
					LifeLeft: 0,
				},
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b"},
			result: Result{
				Id:               "gameId-0",
				Status:           Lose,
				SelectedLetters:  []string{"b"},
				LifeLeft:         0,
				SecretWordLength: 7,
				KnownSecretWord:  "b__b___",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 0,
					LifeLeft: 0,
				},
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o"},
			result: Result{
				Id:               "gameId-0",
				Status:           Lose,
				SelectedLetters:  []string{"b", "o"},
				LifeLeft:         0,
				SecretWordLength: 7,
				KnownSecretWord:  "b__b___",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 0,
					LifeLeft: 0,
				},
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"x", "x"},
			result: Result{
				Id:               "gameId-0",
				Status:           Lose,
				SelectedLetters:  []string{"x"},
				LifeLeft:         0,
				SecretWordLength: 7,
				KnownSecretWord:  "_______",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 0,
					LifeLeft: 0,
				},
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o", "a", "e", "n", "u", "t", "z", "x", "v"},
			result: Result{
				Id:               "gameId-0",
				Status:           Lose,
				SelectedLetters:  []string{"b", "o", "a", "e", "n", "u", "t", "z", "x", "v"},
				LifeLeft:         0,
				SecretWordLength: 7,
				KnownSecretWord:  "b__bea_",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 5,
					LifeLeft: 0,
				},
			},
		},
		{
			inSecretWord: "bigbear",
			inLetters:    []string{"b", "o", "i", "g", "a", "e", "y", "r"},
			result: Result{
				Id:               "gameId-0",
				Status:           Won,
				SelectedLetters:  []string{"b", "o", "i", "g", "a", "e", "y", "r"},
				LifeLeft:         5,
				SecretWordLength: 7,
				KnownSecretWord:  "bigbear",
			},
			timer: Timer{
				Event: GameOver,
				Data: TimerData{
					TimeLeft: 5,
					LifeLeft: 5,
				},
			},
		},
	}
	for _, c := range cases {
		letterCh := make(chan string)
		// hintCh := make(chan string)
		timerCh := make(chan Timer)
		resultByGuess := make(chan Result)
		game := Game{
			Id: "gameId-0",
		}
		resultCh := game.ReactiveHangman(c.inSecretWord, []string{}, letterCh, timerCh, resultByGuess)
		go func() {
			for _, letter := range c.inLetters {
				letterCh <- letter
			}
		}()
		var result Result
		var timer Timer
	GameLoop:
		for {
			select {
			case result = <-resultCh:
				log.Println(result)
				if result.Status == Won || result.Status == Lose {
					break GameLoop
				}
			case timer = <-timerCh:
				log.Println(timer)
			}
		}

		if !reflect.DeepEqual(result, c.result) {
			t.Errorf("Result ReactiveHangman(%q, %q) == %#v, want %#v", c.inSecretWord, c.inLetters, result, c.result)
		}
		if !reflect.DeepEqual(timer, c.timer) {
			t.Errorf("Timer ReactiveHangman(%q, %q) == %#v, want %#v", c.inSecretWord, c.inLetters, timer, c.timer)
		}
	}
}

func TestRandomHintLetters(t *testing.T) {
	cases := []struct {
		inSecretWord string
		want         int
	}{
		{
			inSecretWord: "123",
			want:         1,
		},
		{
			inSecretWord: "1234567891",
			want:         2,
		},
		{
			inSecretWord: "123456789112345",
			want:         3,
		},
	}
	for _, c := range cases {
		game := Game{}
		hintLetters := game.RandomHintLetters(c.inSecretWord)
		if len(hintLetters) != c.want {
			t.Errorf("Result RandomHintLetters(%q) == %#v, want %#v", c.inSecretWord, len(hintLetters), c.want)
		}
	}

}

type MockService struct{}

func (this *MockService) RandomWord() string {
	return "mockword"
}
func (this MockService) RandomHintLetters(secretWord string) []string {
	return []string{"m", "d"}
}

func (this MockService) GetTimerId() string {
	return "mockId"
}

func TestCreate(t *testing.T) {
	want := Result{
		Id:               "gameId-0",
		Status:           InProgress,
		SelectedLetters:  []string{"m", "d"},
		LifeLeft:         7,
		SecretWordLength: 8,
		KnownSecretWord:  "m______d",
	}
	game := Game{
		Id: "gameId-0",
	}
	mockService := MockService{}
	result := <-game.create(&mockService, &mockService, &mockService)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Result create() == %#v, want %#v", result, want)
	}
}
func TestCreateNewGame(t *testing.T) {
	cases := []struct {
		id string
	}{
		{
			id: "gameId-1",
		},
		{
			id: "gameId-2",
		},
		{
			id: "gameId-3",
		},
	}
	for _, c := range cases {
		var got = CreateNewGame()
		if got.Id != c.id {
			t.Errorf("Result CreateNewGame() == %#v, want %#v", got.Id, c.id)
		}
	}
}

func TestGetGameInfo(t *testing.T) {
	cases := []struct {
		id   string
		want Result
	}{
		{
			id:   "gameId-1",
			want: Games["gameId-1"].Result,
		},
		{
			id:   "",
			want: Result{},
		},
	}
	for _, c := range cases {
		var got = GetGameInfo(c.id)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Result GetGameInfo() == %#v, want %#v", got, c.want)
		}
	}
}

func TestGuessLetter(t *testing.T) {
	cases := []struct {
		letter string
		want   Result
	}{
		{
			letter: "x",
			want: Result{
				Id:               "gameId-0",
				Status:           InProgress,
				SelectedLetters:  []string{"m", "d", "x"},
				LifeLeft:         6,
				SecretWordLength: 8,
				KnownSecretWord:  "m______d",
			},
		},
		{
			letter: "o",
			want: Result{
				Id:               "gameId-0",
				Status:           InProgress,
				SelectedLetters:  []string{"m", "d", "x", "o"},
				LifeLeft:         6,
				SecretWordLength: 8,
				KnownSecretWord:  "mo___o_d",
			},
		},
	}
	game := Game{
		Id: "gameId-0",
	}
	Games["gameId-0"] = &game
	mockService := MockService{}
	<-game.create(&mockService, &mockService, &mockService)
	for _, c := range cases {
		got := GuessLetter(game.Id, c.letter)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Result GuessLetter() == %#v, want %#v", got, c.want)
		}
	}

}

func TestSubscibeTimer(t *testing.T) {
	cases := map[int]Timer{}
	cases[1] = Timer{
		Id:    "mockId",
		Event: TimeSpent,
		Data: TimerData{
			TimeLeft: 5,
			LifeLeft: 7,
		},
	}
	cases[2] = Timer{
		Id:    "mockId",
		Event: TimeSpent,
		Data: TimerData{
			TimeLeft: 4,
			LifeLeft: 7,
		},
	}
	cases[6] = Timer{
		Id:    "mockId",
		Event: TimeSpent,
		Data: TimerData{
			TimeLeft: 0,
			LifeLeft: 6,
		},
	}
	cases[6*7] = Timer{
		Id:    "mockId",
		Event: GameOver,
		Data: TimerData{
			TimeLeft: 0,
			LifeLeft: 0,
		},
	}
	// t.Fatal(cases[0])
	game := Game{
		Id: "gameId-0",
	}
	mockService := MockService{}
	timerCh := game.SubscibeTimer()
	<-game.create(&mockService, &mockService, &mockService)
	counter := 0
	for got := range timerCh {
		counter++
		if want, ok := cases[counter]; ok {
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Result SubscibeTimer() == %#v, want %#v", got, want)
			}
		}
		if got.Event == GameOver {
			close(timerCh)
		}
	}

}
