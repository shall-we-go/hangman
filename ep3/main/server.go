package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shall-we-go/hangman/ep3"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hangman", createNewGameHandler).Methods("POST", "GET")
	r.HandleFunc("/hangman/{id}", getGameInfoHandler).Methods("GET")
	r.HandleFunc("/hangman/{id}/timer", getTimerHandler).Methods("GET")
	r.HandleFunc("/hangman/{id}/{letter}", guessLetterHandler).Methods("PUT", "GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createNewGameHandler(w http.ResponseWriter, r *http.Request) {
	result := ep3.CreateNewGame()
	log.Println("createNewGameHandler", result)
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getGameInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result := ep3.GetGameInfo(id)
	log.Println("getGameInfoHandler", result)
	if result.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func guessLetterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	letter := vars["letter"]
	if len(letter) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := ep3.GuessLetter(id, letter)
	log.Println("guessLetterHandler", result)
	if result.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getTimerHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	game := ep3.Games[id]
	if game == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if game.Timer.Event == ep3.GameOver {
		b, err := json.Marshal(game.Timer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	timerCh := game.SubscibeTimer()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	for timer := range timerCh {
		b, err := json.Marshal(timer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\n\n", b)
		flusher.Flush()
		//Stop streaming if Game Over
		if timer.Event == ep3.GameOver {
			return
		}
	}
}
